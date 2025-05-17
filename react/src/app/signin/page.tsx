'use client';
import config from '@/config';
import { useAuth } from '@/contexts/auth_context';
import { useRouter } from 'next/navigation';
import { SignInRequest } from '@/models/api/auth';
import React, { useState } from 'react';
import Link from 'next/link';
import { TextField, Button, Box, Typography, CircularProgress, Alert } from '@mui/material';


export default function SignIn() {
  // auth context
  const context = useAuth();
  const router = useRouter();
  if (!context) {
    throw new Error('Missing auth context');
  }
  // if (context.isAuthenticated) {
  //   router.push('/dashboard');
  //   return null;
  // }

  // state variables for form fields
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');

  // state variables for form submission/ validation
  const [responseText, setResponseText] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState(false);
  const [error, setError] = useState(false);

  // function to handle form submission
  const handleSubmit = async (e: React.FormEvent) => {
    // prevent default form submission
    e.preventDefault();

    // set state variables
    setResponseText(null);
    setLoading(true);
    setSuccess(false);
    setError(false);

    // validate fields
    if (!email || !password) {
      setResponseText('Please fill in all fields');
      setLoading(false);
      setError(true);
      return;
    }

    // create sign in request object
    const signInReq: SignInRequest = {username: '', email, password};
  
    try {
      // make API call to /signin
      const signInRes = await fetch(`${config.servicePath}/signin`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(signInReq),
      });

      // handle response
      if (signInRes.status === 200) {
        const data = await signInRes.json();
        context.login(data.token);
        router.push('/dashboard');
      } else {
        const errorText = await signInRes.text();
        setResponseText(errorText);
        setError(true)
      }

    } catch (err) {
      setResponseText('Something went wrong');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Box maxWidth={400} mx="auto" mt={5}>
      <Typography variant="h4" gutterBottom>Sign In</Typography>
      {error && <Alert severity="error" sx={{ mb: 2 }}>{responseText}</Alert>}
      {success && <Alert severity="success" sx={{ mb: 2 }}>{responseText}</Alert>}
      <form onSubmit={handleSubmit} noValidate>
        <TextField
          label="Email"
          type="email"
          fullWidth
          margin="normal"
          value={email}
          onChange={e => setEmail(e.target.value)}
        />
        <TextField
          label="Password"
          type="password"
          fullWidth
          margin="normal"
          value={password}
          onChange={e => setPassword(e.target.value)}
        />
        <Button
          type="submit"
          variant="contained"
          color="primary"
          fullWidth
          disabled={loading}
          sx={{ mt: 2 }}
        >
          {loading ? <CircularProgress size={24} /> : 'Sign In'}
        </Button>
      </form>
      <Typography variant="body2" align="center" sx={{ mt: 2 }}>
        Don't have an account? <Link href="/signup">Sign Up</Link>
      </Typography>
    </Box>
  );
}
