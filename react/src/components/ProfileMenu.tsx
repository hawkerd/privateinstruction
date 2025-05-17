'use client';

import React from 'react';
import { Menu, MenuItem, ListItemIcon } from '@mui/material';
import SettingsIcon from '@mui/icons-material/Settings';
import LogoutIcon from '@mui/icons-material/Logout';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/contexts/auth_context';

type ProfileMenuProps = {
  anchorEl: null | HTMLElement;
  handleClose: () => void;
};

export default function ProfileMenu({ anchorEl, handleClose }: ProfileMenuProps) {
    // helpful variable
    const open = Boolean(anchorEl);

    // next router
    const router = useRouter();

    // auth context
    const { logout } = useAuth();

    // handle settings click
    const handleSettings = () => {
        handleClose();
        router.push('/settings');
    };

    // handle logout click
    const handleLogout = () => {
        handleClose();
        logout();
        router.push('/');
    };

  return (
    <Menu
      anchorEl={anchorEl}
      open={open}
      onClose={handleClose}
      anchorOrigin={{
        vertical: 'bottom',
        horizontal: 'right',
      }}
      transformOrigin={{
        vertical: 'top',
        horizontal: 'right',
      }}
    >
      <MenuItem onClick={handleSettings}>
        <ListItemIcon>
          <SettingsIcon fontSize="small" />
        </ListItemIcon>
        Settings
      </MenuItem>
      <MenuItem onClick={handleLogout}>
        <ListItemIcon>
          <LogoutIcon fontSize="small" />
        </ListItemIcon>
        Log Out
      </MenuItem>
    </Menu>
  );
}
