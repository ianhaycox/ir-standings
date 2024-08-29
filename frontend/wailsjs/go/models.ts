export namespace live {
	
	export interface PredictedStanding {
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
	}
	export interface Standing {
	    sof_by_car_class: number;
	    car_class_id: number;
	    car_class_name: string;
	    class_leader_laps_complete: number;
	    items: PredictedStanding[];
	}
	export interface PredictedStandings {
	    status: string;
	    track_name: string;
	    count_best_of: number;
	    self_car_class_id: number;
	    car_class_ids: number[];
	    standings: {[key: number]: Standing};
	}

}

export namespace main {
	
	export interface Config {
	    show_topn: number;
	}

}

