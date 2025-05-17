'use client';
import config from '@/config';
import { useAuth } from '@/contexts/auth_context';
import { useRouter } from 'next/navigation';
import { SignInRequest, SignUpRequest } from '@/models/api/auth';
import React, { useState } from 'react';
import { TextField, Button, Box, Typography, CircularProgress, Alert } from '@mui/material';


export default function SignUp() {
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
  const [confirmPassword, setConfirmPassword] = useState('');
  const [username, setUsername] = useState('');

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
    if (!email || !password || !confirmPassword || !username) {
      setResponseText('Please fill in all fields');
      setLoading(false);
      setError(true);
      return;
    }
    if (password !== confirmPassword) {
      setResponseText('Passwords do not match');
      setLoading(false);
      setError(true);
      return;
    }

    // create sign up request object and login request object
    const signUpReq: SignUpRequest = {email, password, username};
    const logInReq: SignInRequest = {username: '', email, password};

    try {
      // make API call to /signup
      const signUpRes = await fetch(`${config.servicePath}/signup`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(signUpReq),
      });

      // handle response
      if (signUpRes.status === 201) {
        setResponseText('Sign up successful. Logging in...');
        setSuccess(true);
      } else {
        const errorText = await signUpRes.text();
        setResponseText(errorText);
        setError(true);
        setLoading(false);
        return;
      }

      // make API call to /login if sign up was successful
      const logInRes = await fetch(`${config.servicePath}/signin`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(logInReq),
      });

      // handle response
      if (logInRes.status === 200) {
        const data = await logInRes.json();
        context.login(data.token);
        router.push('/dashboard');
      } else {
        const errorText = await logInRes.text();
        setResponseText(errorText);
      }

    } catch (err) {
      setResponseText('Something went wrong');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Box maxWidth={400} mx="auto" mt={5}>
      <Typography variant="h4" gutterBottom>Sign Up</Typography>
      {error && <Alert severity="error" sx={{ mb: 2 }}>{responseText}</Alert>}
      {success && <Alert severity="success" sx={{ mb: 2 }}>{responseText}</Alert>}
      <form onSubmit={handleSubmit} noValidate>
        <TextField
          label="Username"
          fullWidth
          margin="normal"
          value={username}
          onChange={e => setUsername(e.target.value)}
        />
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
        <TextField
          label="Confirm Password"
          type="password"
          fullWidth
          margin="normal"
          value={confirmPassword}
          onChange={e => setConfirmPassword(e.target.value)}
        />
        <Button
          type="submit"
          variant="contained"
          color="primary"
          fullWidth
          disabled={loading}
          sx={{ mt: 2 }}
        >
          {loading ? <CircularProgress size={24} /> : 'Sign Up'}
        </Button>
      </form>
    </Box>
  );
}
