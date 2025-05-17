'use client';

import React from 'react';
import { AppBar, Toolbar, Typography, Button, IconButton, Box } from '@mui/material';
import SettingsIcon from '@mui/icons-material/Settings';
import ProfileIcon from '@mui/icons-material/AccountCircle';
import Pets from '@mui/icons-material/Pets';
import Link from 'next/link';
import { useAuth } from '@/contexts/auth_context';

export default function Header() {
    const { token, login, logout, isAuthenticated } = useAuth();

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
                            <IconButton color="inherit">
                                <ProfileIcon />
                            </IconButton>
                            <IconButton color="inherit">
                                <SettingsIcon />
                            </IconButton>
                        </>
                    ) : (
                        <>
                            <Button color="inherit" component={Link} href="/login">
                                Log In
                            </Button>
                            <Button color="inherit" component={Link} href="/signup">
                                Sign Up
                            </Button>
                        </>
                    )}
                </Box>
            </Toolbar>
        </AppBar>
    )
}