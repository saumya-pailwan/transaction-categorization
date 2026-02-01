"use client";

import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { CheckCircle, AlertCircle, HelpCircle } from 'lucide-react';

interface Transaction {
  id: number;
  description: string;
  amount: number;
  category: string;
  merchant_name: string;
  date: string;
}

const TransactionTable = () => {
  const [transactions, setTransactions] = useState<Transaction[]>([]);

  useEffect(() => {
    const fetchTransactions = async () => {
      try {
        const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';
        const response = await axios.get(`${apiUrl}/transactions`);
        setTransactions(response.data || []);
      } catch (error) {
        console.error('Error fetching transactions:', error);
      }
    };
    fetchTransactions();
  }, []);

  return (
    <div className="table-container">
      <table className="transactions-table">
        <thead>
          <tr>
            <th>Date</th>
            <th>Description</th>
            <th>Merchant</th>
            <th>Amount</th>
            <th>Category</th>
            <th>Status</th>
          </tr>
        </thead>
        <tbody>
          {transactions.map((tx) => (
            <tr key={tx.id}>
              <td>{new Date(tx.date).toLocaleDateString()}</td>
              <td className="font-medium">{tx.description}</td>
              <td className="text-secondary">{tx.merchant_name}</td>
              <td className="font-mono amount">${tx.amount.toFixed(2)}</td>
              <td>
                <span className="badge">
                  {tx.category || 'Uncategorized'}
                </span>
              </td>
              <td>
                {tx.category ? (
                  <CheckCircle className="icon-success" />
                ) : (
                  <AlertCircle className="icon-warning" />
                )}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default TransactionTable;