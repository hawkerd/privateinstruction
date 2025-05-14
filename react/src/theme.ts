// src/theme.ts
import { createTheme } from '@mui/material/styles';

const theme = createTheme({
  components: {
    MuiButton: {
      defaultProps: {
        disableRipple: true,
      },
      styleOverrides: {
        root: {
          color: 'black',
          textTransform: 'none',
          boxShadow: 'none',
          '&:hover': {
            backgroundColor: 'transparent',
            boxShadow: 'none',
          },
        },
      },
    },
    MuiIconButton: {
      defaultProps: {
        disableRipple: true,
      },
      styleOverrides: {
        root: {
          color: 'black',
          '&:hover': {
            backgroundColor: 'transparent',
          },
        },
      },
    },
  },
});

export default theme;
