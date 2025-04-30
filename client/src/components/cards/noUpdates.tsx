import React from 'react';
import {
    Box,
    Typography,
    Paper,
    Divider,
    Stack
} from '@mui/material';
import NotificationsOffOutlinedIcon from '@mui/icons-material/NotificationsOffOutlined';

interface NoDataCardProps {
    title?: string;
    description?: string;
    icon?: React.ReactNode;
}

export const NoUpdates: React.FC<NoDataCardProps> = ({
    title = 'No New Updates',
    description = 'You\'re fully up to date. Check back later for any new notifications or system messages.',
    icon
}) => {
    return (
        <Box
            display="flex"
            alignItems="center"
            justifyContent="center"
            height="100vh"
        >
            <Paper
                elevation={0}
                sx={{
                    p: 6,
                    maxWidth: 500,
                    width: '100%',
                    borderRadius: 4,
                    textAlign: 'center',
                    background: '#ffffff',
                    boxShadow: '0 10px 30px rgba(0, 0, 0, 0.08)',
                    border: '1px solid #e0e0e0',
                }}
            >
                <Stack spacing={2} alignItems="center">
                    <Box
                        sx={{
                            backgroundColor: '#1A237E',
                            color: '#ffffff',
                            borderRadius: '50%',
                            width: 64,
                            height: 64,
                            display: 'flex',
                            alignItems: 'center',
                            justifyContent: 'center',
                        }}
                    >
                        {icon || <NotificationsOffOutlinedIcon sx={{ fontSize: 32 }} />}
                    </Box>

                    <Typography variant="h6" fontWeight={600}>
                        {title}
                    </Typography>

                    <Divider sx={{ width: '60%' }} />

                    <Typography variant="body1" color="text.secondary">
                        {description}
                    </Typography>
                </Stack>
            </Paper>
        </Box>
    );
};
