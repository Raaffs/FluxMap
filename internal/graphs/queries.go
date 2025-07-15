package graphs

import "context"

func (g *GraphModel)TaskCompletedByDate(ctx context.Context)(Graphs,error){
	var gs Graphs
	query := `
	SELECT 
	DATE(taskcompleted) AS completed_date
	FROM tasks
	WHERE taskcompleted IS NOT NULL
	GROUP BY completed_date
	HAVING COUNT(*) > 0
	ORDER BY completed_date
	LIMIT 7
	`

	rows,err:=g.DB.Query(ctx,query);if err != nil {
		g.Errorlog.Printf("An error occurred while getting task completed by date: %v\n", err)
		return Graphs{},err
	}
	defer rows.Close()
	for rows.Next() {
		var date string
		var count int
		if err := rows.Scan(&date,&count); err != nil {
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

func (g *GraphModel)TaskApprovedByDate(ctx context.Context)(Graphs,error){
	var gs Graphs
	query := `
	SELECT 
	SELECT 
    DATE(taskapproveddate) AS approved_date
    COUNT(*) 
	FROM tasks 
	WHERE taskcompleteddate IS NOT NULL
	GROUP BY DATE(taskapproveddate)
	ORDER BY approved_date;
	LIMIT 7
	`

	rows,err:=g.DB.Query(ctx,query);if err != nil {
		g.Errorlog.Printf("An error occurred while getting task completed by date: %v\n", err)
		return Graphs{},err
	}
	defer rows.Close()
	for rows.Next() {
		var date string
		var count int
		if err := rows.Scan(&date,&count); err != nil {
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