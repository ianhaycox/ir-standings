export namespace main {
	
	export class PredictedStanding {
	    cust_id: number;
	    driver_name: string;
	    car_number?: string;
	    current_position?: number;
	    predicted_position: number;
	    current_points: number;
	    predicted_points: number;
	    change: number;
	
	    static createFrom(source: any = {}) {
	        return new PredictedStanding(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cust_id = source["cust_id"];
	        this.driver_name = source["driver_name"];
	        this.car_number = source["car_number"];
	        this.current_position = source["current_position"];
	        this.predicted_position = source["predicted_position"];
	        this.current_points = source["current_points"];
	        this.predicted_points = source["predicted_points"];
	        this.change = source["change"];
	    }
	}
	export class PredictedStandings {
	    items: PredictedStanding[];
	
	    static createFrom(source: any = {}) {
	        return new PredictedStandings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.items = this.convertValues(source["items"], PredictedStanding);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

