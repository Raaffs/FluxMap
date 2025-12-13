import * as React from "react";
import dayjs, { Dayjs } from "dayjs";
import { Box, Typography, List, ListItem, ListItemText, Paper } from "@mui/material";
import { DateCalendar, PickersDay, PickersDayProps } from "@mui/x-date-pickers";

interface UpcomingDeadlinesProps {
  deadlines: string[]; // PostgreSQL datetime strings
}

type HighlightedDayProps = PickersDayProps<Dayjs>;
export const UpcomingDeadlines: React.FC<UpcomingDeadlinesProps> = ({ deadlines }) => {
  const deadlineDates: Dayjs[] = deadlines.map((date) => dayjs(date));
  const deadlineDays = deadlineDates.map((d) => d.format("YYYY-MM-DD"));

  function HighlightedDay(props: HighlightedDayProps) {
    const { day, outsideCurrentMonth, ...other } = props;
    const isDeadline = deadlineDays.includes(day.format("YYYY-MM-DD"));

    return (
      <PickersDay
        {...other}
        day={day}
        outsideCurrentMonth={outsideCurrentMonth}
        sx={{
          backgroundColor: isDeadline ? "#ff5252" : undefined,
          color: isDeadline ? "#fff" : undefined,
          borderRadius: "50%",
          "&:hover": {
            backgroundColor: isDeadline ? "#ff1744" : undefined
          }
        }}
      />
    );
  }

  const upcomingList = [...deadlineDates]
    .filter((d) => d.isAfter(dayjs().subtract(1, "day")))
    .sort((a, b) => a.valueOf() - b.valueOf())
    .slice(0, 5);

  return (
    <Paper sx={{ p: 2, width: "100%" }}>
      <Typography variant="h6" gutterBottom>
        Upcoming Deadlines
      </Typography>

      <Box display="flex" gap={3}>
        <DateCalendar slots={{ day: HighlightedDay }} views={["day"]} />
        <List dense>
          {upcomingList.map((date, idx) => (
            <ListItem key={idx}>
              <ListItemText
                primary={date.format("DD MMM YYYY")}
                secondary={date.format("hh:mm A")}
              />
            </ListItem>
          ))}
        </List>
      </Box>
    </Paper>
  );
};
