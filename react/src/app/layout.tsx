// src/app/layout.tsx
'use client';
import { AuthProvider } from '@/contexts/auth_context';
import { ThemeProvider } from '@mui/material/styles';
import { UserProvider } from '@/contexts/user_context';
import CssBaseline from '@mui/material/CssBaseline';
import theme from '@/theme';
import Header from '@/components/Header';
import { usePathname } from 'next/navigation';

export default function RootLayout({ children }: { children: React.ReactNode }) {
  const pathName = usePathname();
  const noHeaderPaths = ['/signin', '/signup'];
  const showHeader = !noHeaderPaths.includes(pathName);



  return (
    <html lang="en">
      <body style={{ margin: 0 }}>
        <AuthProvider>
          <UserProvider>
            <ThemeProvider theme={theme}>
              <CssBaseline />
              {showHeader && <Header />}
              {children}
            </ThemeProvider>
          </UserProvider>
        </AuthProvider>
      </body>
    </html>
  );
}
