"use client";

import React, { useCallback, useEffect, useState } from 'react';
import { usePlaidLink } from 'react-plaid-link';
import axios from 'axios';

const PlaidLink = () => {
  const [token, setToken] = useState<string | null>(null);

  useEffect(() => {
    const createLinkToken = async () => {
      try {
        const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';
        const response = await axios.post(`${apiUrl}/plaid/create_link_token`);
        setToken(response.data.link_token);
      } catch (error) {
        console.error('Error creating link token:', error);
      }
    };
    createLinkToken();
  }, []);

  const onSuccess = useCallback(async (publicToken: string) => {
    try {
      const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';
      await axios.post(`${apiUrl}/plaid/exchange_public_token`, {
        public_token: publicToken,
      });
      // Trigger a refresh or notify parent
      window.location.reload();
    } catch (error) {
      console.error('Error exchanging public token:', error);
    }
  }, []);

  const config: Parameters<typeof usePlaidLink>[0] = {
    token,
    onSuccess,
  };

  const { open, ready } = usePlaidLink(config);

  return (
    <button
      onClick={() => open()}
      disabled={!ready}
      className="btn-primary"
    >
      Connect Bank Account
    </button>
  );
};

export default PlaidLink;
