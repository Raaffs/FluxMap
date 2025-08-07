package graphs

import "context"

func (g *GraphModel)TaskCompletedByDate(ctx context.Context, assignedUsername string)(Graphs,error){
	var gs Graphs
	query := `
	SELECT 
		COUNT(*) AS task_count, 
       	DATE(taskcompleteddate) AS completed_date
		FROM tasks
		WHERE assignedUsername = $1
		AND taskcompleteddate IS NOT NULL
		AND taskcompleteddate >= CURRENT_DATE - INTERVAL '100 days'
		GROUP BY completed_date
		ORDER BY completed_date;
	`

	rows,err:=g.DB.Query(ctx,query,assignedUsername);if err != nil {
		g.Errorlog.Printf("An error occurred while getting task completed by date: %v\n", err)
		return Graphs{},err
	}
	defer rows.Close()
	for rows.Next() {
		var count int
		var date any
		if err := rows.Scan(&count,&date); err != nil {
			g.Errorlog.Printf("An error occurred while scanning task completed by date: %v\n", err)
			return Graphs{},err
		}
		gs.XAxis = append(gs.XAxis, date)
		gs.YAxis = append(gs.YAxis, count)	
	}
	if rows.Err() != nil {
		g.Errorlog.Printf("An error occurred while getting task completed by date: %v\n", rows.Err())
		return Graphs{},rows.Err()
	}

	return gs,nil
}

func (g *GraphModel)TaskApprovedByDate(ctx context.Context, assginedUsername string)(Graphs,error){
	var gs Graphs
	query := `
	SELECT     
		COUNT(*) AS task_count,
    	DATE(taskapproveddate) AS approved_date
		FROM tasks 
		WHERE assignedUsername= $1 
		AND taskapproveddate IS NOT NULL
		AND taskapproveddate >= CURRENT_DATE - INTERVAL '100 days'
		GROUP BY approved_date
		ORDER BY approved_date
	`

	rows,err:=g.DB.Query(ctx,query,assginedUsername);if err != nil {
		g.Errorlog.Printf("An error occurred while getting task completed by date: %v\n", err)
		return Graphs{},err
	}
	defer rows.Close()
	for rows.Next() {
		var date any
		var count int
		if err := rows.Scan(&count,&date); err != nil {
			g.Errorlog.Printf("An error occurred while scanning task completed by date: %v\n", err)
			return Graphs{},err
		}
		gs.XAxis = append(gs.XAxis, date)
		gs.YAxis = append(gs.YAxis, count)	
	}
	if rows.Err() != nil {
		g.Errorlog.Printf("An error occurred while getting task completed by date: %v\n", rows.Err())
		return Graphs{},rows.Err()
	}

	return gs,nil
}
