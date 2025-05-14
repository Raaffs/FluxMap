import * as React from 'react';
import Snackbar from '@mui/material/Snackbar';
import Alert from '@mui/material/Alert';

interface PopUpProps {
    Message?: string;
    Error?: string;
    onClose?: () => void;
}

const PopUp: React.FC<PopUpProps> = ({ Message, Error, onClose }) => {
    const [open, setOpen] = React.useState<boolean>(true);

    const handleClose = (
        event?: React.SyntheticEvent | Event,
        reason?: string
    ) => {
        if (reason === 'clickaway') {
            return;
        }
        setOpen(false);
        if (onClose) {
            onClose(); // Notify parent that Snackbar has closed
        }
    };

    React.useEffect(() => {
        if (Message || Error) {
            setOpen(true); // Open Snackbar for new messages or errors
        }
    }, [Message, Error]);

    return (
        <div>
            <Snackbar
                open={open}
                autoHideDuration={5000}
                onClose={handleClose}
                anchorOrigin={{ vertical: 'top', horizontal: 'center' }} // Set position at the top center
            >
                <Alert
                    onClose={handleClose}
                    severity={Error ? 'error' : 'success'} // Set severity based on Error prop
                    variant="filled"
                    sx={{
                        width: '200%',
                        fontSize: '1rem',
                    }}
                >
                    {Error ? Error : Message} {/* Show Error if present, otherwise show success Message */}
                </Alert>
            </Snackbar>
        </div>
    );
};

export default PopUp;
