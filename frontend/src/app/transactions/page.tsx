"use client";
import { useEffect, useState } from "react";
import { getTransactions } from "@/lib/api";
import TransactionTable from "@/components/TransactionTable";

export default function TransactionsPage() {
  const [tx, setTx] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function fetchData() {
      const data = await getTransactions();
      setTx(data || []);
      setLoading(false);
    }
    fetchData();
  }, []);

  if (loading) return <p>Loading...</p>;

  return (
    <div>
      <h2 className="text-2xl font-semibold mb-4">Transactions</h2>
      <TransactionTable transactions={tx} />
    </div>
  );
}