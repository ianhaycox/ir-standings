export namespace live {
	
	export class PredictedStanding {
	    driving: boolean;
	    cust_id: number;
	    driver_name: string;
	    car_number: string;
	    current_position: number;
	    predicted_position: number;
	    current_points: number;
	    predicted_points: number;
	    change: number;
	    car_names: string[];
	
	    static createFrom(source: any = {}) {
	        return new PredictedStanding(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.driving = source["driving"];
	        this.cust_id = source["cust_id"];
	        this.driver_name = source["driver_name"];
	        this.car_number = source["car_number"];
	        this.current_position = source["current_position"];
	        this.predicted_position = source["predicted_position"];
	        this.current_points = source["current_points"];
	        this.predicted_points = source["predicted_points"];
	        this.change = source["change"];
	        this.car_names = source["car_names"];
	    }
	}
	export class Standing {
	    sof_by_car_class: number;
	    car_class_id: number;
	    car_class_name: string;
	    class_leader_laps_complete: number;
	    items: PredictedStanding[];
	
	    static createFrom(source: any = {}) {
	        return new Standing(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sof_by_car_class = source["sof_by_car_class"];
	        this.car_class_id = source["car_class_id"];
	        this.car_class_name = source["car_class_name"];
	        this.class_leader_laps_complete = source["class_leader_laps_complete"];
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
	export class PredictedStandings {
	    status: string;
	    track_name: string;
	    count_best_of: number;
	    self_car_class_id: number;
	    car_class_ids: number[];
	    standings: {[key: number]: Standing};
	
	    static createFrom(source: any = {}) {
	        return new PredictedStandings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.track_name = source["track_name"];
	        this.count_best_of = source["count_best_of"];
	        this.self_car_class_id = source["self_car_class_id"];
	        this.car_class_ids = source["car_class_ids"];
	        this.standings = this.convertValues(source["standings"], Standing, true);
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

