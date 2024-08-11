/*
MIT License

Copyright (c) 2021-2022 L. E. Spalt

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

#pragma once

#include <assert.h>
#include "Overlay.h"
#include "Config.h"
#include "OverlayDebug.h"
#include "live.h"

class OverlayStandings : public Overlay
{
public:

    const float DefaultFontSize = 15;

    enum class Columns { PREDICTED_STANDING, CAR_NUMBER, NAME, CURRENT_STANDING, POINTS, CHANGE };

    OverlayStandings(const int selectedClassID)
        : Overlay("OverlayStandings" + std::to_string(selectedClassID))
    {
        m_selectedClassID = selectedClassID;
        m_name = "OverlayStandings" + std::to_string(selectedClassID);
    }

private:
    int m_selectedClassID;

protected:

    virtual void onEnable()
    {
        onConfigChanged();  // trigger font load
    }

    virtual void onDisable()
    {
        m_text.reset();
    }

    virtual void onConfigChanged()
    {
        m_text.reset( m_dwriteFactory.Get() );

        const std::string font = g_cfg.getString( m_name, "font", "Microsoft YaHei UI" );
        const float fontSize = g_cfg.getFloat( m_name, "font_size", DefaultFontSize );
        const int fontWeight = g_cfg.getInt( m_name, "font_weight", 500 );
        HRCHECK(m_dwriteFactory->CreateTextFormat( toWide(font).c_str(), NULL, (DWRITE_FONT_WEIGHT)fontWeight, DWRITE_FONT_STYLE_NORMAL, DWRITE_FONT_STRETCH_NORMAL, fontSize, L"en-us", &m_textFormat ));
        m_textFormat->SetParagraphAlignment( DWRITE_PARAGRAPH_ALIGNMENT_CENTER );
        m_textFormat->SetWordWrapping( DWRITE_WORD_WRAPPING_NO_WRAP );

        HRCHECK(m_dwriteFactory->CreateTextFormat( toWide(font).c_str(), NULL, (DWRITE_FONT_WEIGHT)fontWeight, DWRITE_FONT_STYLE_NORMAL, DWRITE_FONT_STRETCH_NORMAL, fontSize*0.8f, L"en-us", &m_textFormatSmall ));
        m_textFormatSmall->SetParagraphAlignment( DWRITE_PARAGRAPH_ALIGNMENT_CENTER );
        m_textFormatSmall->SetWordWrapping( DWRITE_WORD_WRAPPING_NO_WRAP );

        // Determine widths of text columns
        m_columns.reset();
        m_columns.add((int)Columns::PREDICTED_STANDING, computeTextExtent(L"E9999", m_dwriteFactory.Get(), m_textFormat.Get()).x, fontSize / 2);
        m_columns.add((int)Columns::CAR_NUMBER, computeTextExtent(L"#9999", m_dwriteFactory.Get(), m_textFormat.Get()).x, fontSize / 2);
        m_columns.add((int)Columns::NAME, 0, fontSize / 2);
        m_columns.add((int)Columns::CURRENT_STANDING,  computeTextExtent(L"C9999", m_dwriteFactory.Get(), m_textFormat.Get()).x, fontSize / 2);
        m_columns.add((int)Columns::POINTS, computeTextExtent(L"9999", m_dwriteFactory.Get(), m_textFormat.Get()).x, fontSize / 2);
        m_columns.add((int)Columns::CHANGE,            computeTextExtent(L"999", m_dwriteFactory.Get(), m_textFormat.Get()).x, fontSize / 2);
    }

    virtual void onUpdate()
    {
        struct CarInfo {
            int     carIdx = 0;
            int     position = 0;
            int     change = 0;
            int     lapsComplete = 0;
        };

        std::vector<CarInfo> carInfo;
        std::map<int,std::string> carNumbers;
        int classLeaderLapsComplete = 0;

        carInfo.reserve( IR_MAX_CARS );

        // Init array
        for( int i=0; i<IR_MAX_CARS; ++i )
        {
            const Car& car = ir_session.cars[i];

            if( car.isPaceCar || car.isSpectator || car.userName.empty() || car.carClassID != m_selectedClassID)
                continue;

            CarInfo ci;
            ci.carIdx       = i;
            ci.position     = ir_getPosition(i);
            ci.lapsComplete = ir_CarIdxLapCompleted.getInt(i);
            carNumbers[car.custID] = car.carNumberStr;

            carInfo.push_back( ci );
        }

        struct LiveResults lr;
        lr.seriesID = ir_session.seriesId;
        lr.sessionID = ir_session.sessionId;
        lr.subsessionID = ir_session.subsessionId;
        lr.track = ir_session.trackName;
        lr.countBestOf = g_cfg.getInt(m_name, "count_best_of", 10);
        lr.carClassID = m_selectedClassID;
        lr.topN = g_cfg.getInt(m_name, "top_n", 20);

        for (int i = 0; i<carInfo.size(); ++i) {
            struct CurrentPosition cp;

            const Car&  car = ir_session.cars[carInfo[i].carIdx];

            cp.carID = car.carID;
            cp.finishPositionInClass = carInfo[i].position;
            cp.lapsComplete = carInfo[i].lapsComplete;
            cp.custID = car.custID;

            if (cp.lapsComplete > classLeaderLapsComplete) {
                classLeaderLapsComplete = cp.lapsComplete;
            }

            lr.positions.push_back(cp);
        }

        std::string fn = g_cfg.getString("General", "filename", "285-results.json");

        std::vector<PredictedStanding> predictedStandings  = LatestStandings(fn, lr);

        const float  fontSize           = g_cfg.getFloat( m_name, "font_size", DefaultFontSize );
        const float  lineSpacing        = g_cfg.getFloat( m_name, "line_spacing", 8 );
        const float  lineHeight         = fontSize + lineSpacing;
        const float4 selfCol            = g_cfg.getFloat4( m_name, "self_col", float4(0.94f,0.67f,0.13f,1) );
        const float4 buddyCol           = g_cfg.getFloat4( m_name, "buddy_col", float4(0.2f,0.75f,0,1) );
        const float4 flaggedCol         = g_cfg.getFloat4( m_name, "flagged_col", float4(0.68f,0.42f,0.2f,1) );
        const float4 otherCarCol        = g_cfg.getFloat4( m_name, "other_car_col", float4(1,1,1,0.9f) );
        const float4 headerCol          = g_cfg.getFloat4( m_name, "header_col", float4(0.7f,0.7f,0.7f,0.9f) );
        const float4 carNumberTextCol   = g_cfg.getFloat4( m_name, "car_number_text_col", float4(0,0,0,0.9f) );
        const float4 carNumberBgCol     = g_cfg.getFloat4(m_name,  "car_number_background_col", float4(0, 0, 0, 0.7f));
        const float4 alternateLineBgCol = g_cfg.getFloat4( m_name, "alternate_line_background_col", float4(0.5f,0.5f,0.5f,0.1f) );
        const float4 changePlusBgCol    = g_cfg.getFloat4( m_name, "change_plus_background_col", float4(0.19f,0.58f,0.12f,1) );
        const float4 changeMinusBgCol   = g_cfg.getFloat4( m_name, "change_minus_background_col", float4(0.65f,0.11f,0.11f,1) );
        const bool   imperial           = ir_DisplayUnits.getInt() == 0;

        const float xoff = 10.0f;
        const float yoff = 10;
        m_columns.layout( (float)m_width - 2*xoff );
        float y = yoff + lineHeight/2;
        const float ybottom = m_height - lineHeight * 1.5f;

        const ColumnLayout::Column* clm = nullptr;
        wchar_t s[512];
        std::string str;
        D2D1_RECT_F r = {};
        D2D1_ROUNDED_RECT rr = {};

        m_renderTarget->BeginDraw();
        m_brush->SetColor( headerCol );

        // Headers
        clm = m_columns.get((int)Columns::PREDICTED_STANDING);
        swprintf(s, _countof(s), L"Pos.");
        m_text.render(m_renderTarget.Get(), s, m_textFormat.Get(), xoff + clm->textL, xoff + clm->textR, y, m_brush.Get(), DWRITE_TEXT_ALIGNMENT_LEADING);

        clm = m_columns.get((int)Columns::CAR_NUMBER);
        swprintf(s, _countof(s), L"No.");
        m_text.render(m_renderTarget.Get(), s, m_textFormat.Get(), xoff + clm->textL, xoff + clm->textR, y, m_brush.Get(), DWRITE_TEXT_ALIGNMENT_LEADING);

        clm = m_columns.get((int)Columns::NAME);
        swprintf(s, _countof(s), L"Driver");
        m_text.render(m_renderTarget.Get(), s, m_textFormat.Get(), xoff + clm->textL, xoff + clm->textR, y, m_brush.Get(), DWRITE_TEXT_ALIGNMENT_LEADING);

        clm = m_columns.get((int)Columns::CURRENT_STANDING);
        swprintf(s, _countof(s), L"Prev.");
        m_text.render(m_renderTarget.Get(), s, m_textFormat.Get(), xoff + clm->textL, xoff + clm->textR, y, m_brush.Get(), DWRITE_TEXT_ALIGNMENT_LEADING);

        clm = m_columns.get( (int)Columns::POINTS );
        swprintf( s, _countof(s), L"Pts." );
        m_text.render( m_renderTarget.Get(), s, m_textFormat.Get(), xoff+clm->textL, xoff+clm->textR, y, m_brush.Get(), DWRITE_TEXT_ALIGNMENT_TRAILING);

        clm = m_columns.get((int)Columns::CHANGE);
        swprintf(s, _countof(s), L"+/-");
        m_text.render(m_renderTarget.Get(), s, m_textFormat.Get(), xoff + clm->textL, xoff + clm->textR, y, m_brush.Get(), DWRITE_TEXT_ALIGNMENT_CENTER);

        // Content
        for( int i=0; i<predictedStandings.size(); ++i )
        {
            float4 textCol = otherCarCol;
            const bool notPresent = carNumbers[predictedStandings[i].custID].empty();
            if (notPresent) {
                textCol.a *= 0.5f;
            }

            y = 2*yoff + lineHeight/2 + (i+1)*lineHeight;

            if( y+lineHeight/2 > ybottom )
                break;

            // Alternating line backgrounds
            if( i & 1 && alternateLineBgCol.a > 0 )
            {
                D2D1_RECT_F r = { 0, y-lineHeight/2, (float)m_width,  y+lineHeight/2 };
                m_brush->SetColor( alternateLineBgCol );
                m_renderTarget->FillRectangle( &r, m_brush.Get() );
            }

            {
                clm = m_columns.get((int)Columns::PREDICTED_STANDING);
                m_brush->SetColor(textCol);
                swprintf(s, _countof(s), L"P%d", predictedStandings[i].predictedPosition);
                m_text.render(m_renderTarget.Get(), s, m_textFormat.Get(), xoff + clm->textL, xoff + clm->textR, y, m_brush.Get(), DWRITE_TEXT_ALIGNMENT_LEADING);
            }

            // Car number
            {
                clm = m_columns.get( (int)Columns::CAR_NUMBER );
                if (carNumbers[predictedStandings[i].custID].empty()) {
                    m_brush->SetColor( textCol );
                    swprintf( s, _countof(s), L"#%S",  "---" );
                } else {
                    swprintf( s, _countof(s), L"#%S",  carNumbers[predictedStandings[i].custID].c_str() );
                    r = { xoff+clm->textL, y-lineHeight/2, xoff+clm->textR, y+lineHeight/2 };
                    rr.rect = { r.left-2, r.top+1, r.right+2, r.bottom-1 };
                    rr.radiusX = 3;
                    rr.radiusY = 3;
                    m_brush->SetColor(carNumberBgCol);
                    m_renderTarget->FillRoundedRectangle( &rr, m_brush.Get() );
                    m_brush->SetColor( carNumberTextCol );
                }
                m_text.render( m_renderTarget.Get(), s, m_textFormat.Get(), xoff+clm->textL, xoff+clm->textR, y, m_brush.Get(), DWRITE_TEXT_ALIGNMENT_LEADING);
            }

            // Name
            {
                clm = m_columns.get( (int)Columns::NAME );
                m_brush->SetColor( textCol );
                swprintf( s, _countof(s), L"%S", predictedStandings[i].driverName.c_str() );
                m_text.render( m_renderTarget.Get(), s, m_textFormat.Get(), xoff+clm->textL, xoff+clm->textR, y, m_brush.Get(), DWRITE_TEXT_ALIGNMENT_LEADING );
            }

            {
                clm = m_columns.get((int)Columns::CURRENT_STANDING);
                float4 dim = textCol;
                dim.a *= 0.5f;
                m_brush->SetColor(dim);
                swprintf(s, _countof(s), L"P%d", predictedStandings[i].currentPosition);
                m_text.render(m_renderTarget.Get(), s, m_textFormat.Get(), xoff + clm->textL, xoff + clm->textR, y, m_brush.Get(), DWRITE_TEXT_ALIGNMENT_LEADING);
            }

            {
                clm = m_columns.get((int)Columns::POINTS);
                m_brush->SetColor(textCol);
                swprintf(s, _countof(s), L"%d", predictedStandings[i].predictedPoints);
                m_text.render(m_renderTarget.Get(), s, m_textFormat.Get(), xoff + clm->textL, xoff + clm->textR, y, m_brush.Get(), DWRITE_TEXT_ALIGNMENT_TRAILING);
            }

            {
                clm = m_columns.get((int)Columns::CHANGE);
                r = { xoff+clm->textL, y-lineHeight/2, xoff+clm->textR, y+lineHeight/2 };
                rr.rect = { r.left-2, r.top+1, r.right+2, r.bottom-1 };
                rr.radiusX = 3;
                rr.radiusY = 3;

                if (predictedStandings[i].change > 0) {
                    swprintf(s, _countof(s), L"^ %d", predictedStandings[i].change);
                    m_brush->SetColor(changePlusBgCol);
                    m_renderTarget->FillRoundedRectangle( &rr, m_brush.Get() );
                }

                if (predictedStandings[i].change < 0) {
                    swprintf(s, _countof(s), L"v %d", predictedStandings[i].change*-1);
                    m_brush->SetColor(changeMinusBgCol);
                    m_renderTarget->FillRoundedRectangle( &rr, m_brush.Get() );
                }

                if (predictedStandings[i].change == 0) {
                    swprintf(s, _countof(s), L"%S", "---");
                    float4 dim = textCol;
                    dim.a *= 0.5f;
                    m_brush->SetColor(dim);
                    m_renderTarget->FillRoundedRectangle( &rr, m_brush.Get() );
                }

                m_brush->SetColor( textCol );
                m_text.render(m_renderTarget.Get(), s, m_textFormat.Get(), xoff + clm->textL, xoff + clm->textR, y, m_brush.Get(), DWRITE_TEXT_ALIGNMENT_CENTER);
            }
        }

        // Footer
        {
            m_brush->SetColor(float4(1,1,1,0.4f));
            m_renderTarget->DrawLine( float2(0,ybottom),float2((float)m_width,ybottom),m_brush.Get() );
            swprintf( s, _countof(s), L"%S  Best Of: %d  SoF: %d  Laps: %d", ir_session.trackName.c_str(), lr.countBestOf, ir_session.sofByCarClass[m_selectedClassID], classLeaderLapsComplete);
            y = m_height - (m_height-ybottom)/2;
            m_brush->SetColor( headerCol );
            m_text.render( m_renderTarget.Get(), s, m_textFormat.Get(), xoff, (float)m_width-2*xoff, y, m_brush.Get(), DWRITE_TEXT_ALIGNMENT_CENTER );
        }

        m_renderTarget->EndDraw();
    }

    virtual bool canEnableWhileNotDriving() const
    {
        return true;
    }

protected:

    Microsoft::WRL::ComPtr<IDWriteTextFormat>  m_textFormat;
    Microsoft::WRL::ComPtr<IDWriteTextFormat>  m_textFormatSmall;

    ColumnLayout m_columns;
    TextCache    m_text;
};
