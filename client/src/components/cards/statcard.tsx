import { Box, Typography, useTheme, ButtonBase } from "@mui/material";
import { tokens } from "../../theme";
import { useNavigate } from "react-router-dom";

interface StatCardProps {
  title: string;
  subtitle: string;
  icon: React.ReactNode;
  fontColor: string;
  path: string;
}

const StatCard = ({
  title,
  subtitle,
  icon,
  fontColor,
  path,
}: StatCardProps) => {
  const theme = useTheme();
  const colors = tokens(theme.palette.mode);
  const navigate = useNavigate();

  const handleClick = () => {
    navigate(path);
  };

  return (
    <ButtonBase
      onClick={handleClick}
      sx={{
        width: "100%",
        m: "0 30px",
        borderRadius: 2,
        display: "block",
        textAlign: "center",
        padding: 2,
        transition: "background-color 0.2s",
        "&:hover": {
          backgroundColor: colors.primary[400],
        },
      }}
    >
      <Box display="flex" flexDirection="column" alignItems="center">
        <Box display="flex" justifyContent="center" alignItems="center" gap="12px">
          <Typography variant="h4" fontWeight="bold" sx={{ color: fontColor }}>
            {title}
          </Typography>
          <Box display="flex" justifyContent="center" alignItems="center">
            {icon}
          </Box>
        </Box>

        <Box mt="2px">
          <Typography variant="h5" sx={{ color: fontColor }}>
            {subtitle}
          </Typography>
          <Typography
            variant="h5"
            fontStyle="italic"
            sx={{ color: colors.greenAccent[600] }}
          />
        </Box>
      </Box>
    </ButtonBase>
  );
};

export default StatCard;
