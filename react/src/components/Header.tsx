'use client';

import React from 'react';
import { AppBar, Toolbar, Typography, Button, IconButton, Box } from '@mui/material';
import SettingsIcon from '@mui/icons-material/Settings';
import ProfileIcon from '@mui/icons-material/AccountCircle';
import Pets from '@mui/icons-material/Pets';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/contexts/auth_context';
import { useUser } from '@/contexts/user_context';
import ProfileMenu from '@/components/ProfileMenu';

export default function Header() {
    const { token, login, logout, isAuthenticated } = useAuth();
    const { user, setUser, clearUser } = useUser();

    // state variable for the profile menu
    const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);

    // functions for the profile menu
    const handleClose = () => {
        setAnchorEl(null);
    }
    const handleClick = (event: React.MouseEvent<HTMLElement>) => {
        setAnchorEl(event.currentTarget);
    }

    return (
        <AppBar
            variant='outlined'
            position="static"
            sx={{ boxShadow: 'none', margin: 0, padding: 0, border: 'none' }}
        >
            <Toolbar 
                sx={{backgroundColor: 'white', display: 'flex', justifyContent: 'space-between'}}
            >
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 3 }}>
                    <IconButton 
                        sx={{ color: 'black', scale: 1.5 }}
                        component={Link}
                        href="/"
                    >
                        <Pets />
                    </IconButton>
                    <Box sx={{ display: 'flex', alignItems: 'center' }}>
                        <Button
                            sx={{ color: 'black' }}
                            component={Link}
                            href="/about">
                            About
                        </Button>
                        <Button
                            sx={{ color: 'black' }}
                            component={Link}
                            href="/plans">
                            Plans
                        </Button>
                        <Button
                            sx={{ color: 'black' }}
                            component={Link}
                            href="/support">
                            Support
                        </Button>
                    </Box>
                </Box>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                    {isAuthenticated ? (
                        <>
                            {user && <Typography variant="h6" sx={{ color: 'black' }}>{user.username}</Typography>}
                            <IconButton
                                sx={{ color: 'black' }}
                                onClick={handleClick}>
                                <ProfileIcon />
                            </IconButton>
                            <ProfileMenu
                                anchorEl={anchorEl}
                                handleClose={handleClose}
                            />
                        </>

                    ) : (
                        <Button
                            sx={{ color: 'black' }}
                            component={Link}
                            href="/signin">
                            Sign In
                        </Button>
                    )}
                </Box>
            </Toolbar>
        </AppBar>
    )
}