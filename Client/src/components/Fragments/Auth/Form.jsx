"use client";

import { Button, Label, TextInput } from "flowbite-react";
import axios from 'axios';
import React, { useState } from 'react';
import { jwtDecode } from 'jwt-decode';

export function LoginForm() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');

  const handleLogin = async (e) => {
    e.preventDefault();
    try {
      const response = await axios.post('http://localhost:8080/login', {
        email, 
        password
      });
      localStorage.setItem('token', response.data.token);
      const decoded = jwtDecode(response.data.token);
      console.log('Role:', decoded.role); // Log role ke console
      // Redirect ke dashboard
      window.location.href = '/dashboard';
    } catch (error) {
      alert('Login gagal: ' + error.response.data.error);
    }
  };

  return (
    <form onSubmit={handleLogin}>
      <div>
        <div className="mb-2 block">
          <Label className="text-white" htmlFor="email" value="Email" />
        </div>
        <TextInput
          id="email"
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          placeholder="contoh@mail.com"
          required
        />
      </div>
      <div>
        <div className="mb-2 block">
          <Label className="text-white" htmlFor="password" value="Password" />
        </div>
        <TextInput
          id="password"
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          placeholder="••••••••"
          required
        />
      </div>
      <Button color="warning" type="submit">
        Kirim
      </Button>
    </form>
  );
}

export function RegisterForm() {
  const handleRegister = async (e) => {
    e.preventDefault();
    const userData = {
      username: document.getElementById('username').value,
      email: document.getElementById('email').value,
      password: document.getElementById('passsword').value,
      role: document.getElementById('role').value
    };

    try {
      const response = await axios.post('http://localhost:8080/register', userData);
      alert('Registrasi berhasil!');
    } catch (error) {
      alert('Registrasi gagal: ' + error.response.data.error);
    }
  };

  return (
    <form onSubmit={handleRegister} className="flex w-full flex-col gap-4">
      <div>
        <div className="mb-2 block">
          <Label className="text-white" htmlFor="username" value="Username" />
        </div>
        <TextInput id="username" type="text" required />
      </div>
      <div>
        <div className="mb-2 block">
          <Label className="text-white" htmlFor="email" value="Alamat Email" />
        </div>
        <TextInput
          id="email"
          type="email"
          placeholder="contoh@mail.com"
          required
        />
      </div>
      <div>
        <div className="mb-2 block">
          <Label className="text-white" htmlFor="passsword" value="Password" />
        </div>
        <TextInput
          id="passsword"
          type="password"
          placeholder="••••••••"
          required
        />
      </div>
      <div>
        <div className="mb-2 block">
          <Label className="text-white" htmlFor="role" value="Role" />
        </div>
        <select id="role" className="form-select block w-full px-3 py-1.5 text-base font-normal text-gray-700 bg-white bg-clip-padding bg-no-repeat border border-solid border-gray-300 rounded transition ease-in-out m-0 focus:text-gray-700 focus:bg-white focus:border-blue-600 focus:outline-none" required>
          <option value="user">User</option>
          <option value="admin">Admin</option>
        </select>
      </div>
      <Button color="warning" type="submit">
        Kirim
      </Button>
    </form>
  );
}