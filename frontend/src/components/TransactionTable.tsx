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
        const response = await axios.get('http://localhost:8080/api/transactions');
        setTransactions(response.data || []);
      } catch (error) {
        console.error('Error fetching transactions:', error);
      }
    };
    fetchTransactions();
  }, []);

  return (
    <div className="overflow-x-auto">
      <table className="min-w-full bg-white border border-gray-200">
        <thead>
          <tr className="bg-gray-50">
            <th className="py-2 px-4 border-b text-left">Date</th>
            <th className="py-2 px-4 border-b text-left">Description</th>
            <th className="py-2 px-4 border-b text-left">Merchant</th>
            <th className="py-2 px-4 border-b text-left">Amount</th>
            <th className="py-2 px-4 border-b text-left">Category</th>
            <th className="py-2 px-4 border-b text-left">Status</th>
          </tr>
        </thead>
        <tbody>
          {transactions.map((tx) => (
            <tr key={tx.id} className="hover:bg-gray-50">
              <td className="py-2 px-4 border-b">{new Date(tx.date).toLocaleDateString()}</td>
              <td className="py-2 px-4 border-b">{tx.description}</td>
              <td className="py-2 px-4 border-b">{tx.merchant_name}</td>
              <td className="py-2 px-4 border-b font-mono">${tx.amount.toFixed(2)}</td>
              <td className="py-2 px-4 border-b">
                <span className="inline-block bg-gray-100 rounded-full px-3 py-1 text-sm font-semibold text-gray-700">
                  {tx.category || 'Uncategorized'}
                </span>
              </td>
              <td className="py-2 px-4 border-b">
                {tx.category ? (
                  <CheckCircle className="text-green-500 w-5 h-5" />
                ) : (
                  <AlertCircle className="text-yellow-500 w-5 h-5" />
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