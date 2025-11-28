"use client";

import React, { useCallback, useEffect, useState } from 'react';
import { usePlaidLink } from 'react-plaid-link';
import axios from 'axios';

const PlaidLink = () => {
  const [token, setToken] = useState<string | null>(null);

  useEffect(() => {
    const createLinkToken = async () => {
      try {
        const response = await axios.post('http://localhost:8080/api/plaid/create_link_token');
        setToken(response.data.link_token);
      } catch (error) {
        console.error('Error creating link token:', error);
      }
    };
    createLinkToken();
  }, []);

  const onSuccess = useCallback(async (publicToken: string) => {
    try {
      await axios.post('http://localhost:8080/api/plaid/exchange_public_token', {
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
      className="bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
    >
      Connect Bank Account
    </button>
  );
};

export default PlaidLink;
