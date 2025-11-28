import PlaidLink from '@/components/PlaidLink';
import TransactionTable from '@/components/TransactionTable';

export default function Home() {
  return (
    <main className="flex min-h-screen flex-col items-center p-24 bg-gray-50">
      <div className="z-10 max-w-5xl w-full items-center justify-between font-mono text-sm lg:flex mb-8">
        <h1 className="text-4xl font-bold text-gray-900">Autonomous Transaction Agent</h1>
        <PlaidLink />
      </div>

      <div className="w-full max-w-5xl">
        <div className="bg-white shadow-md rounded-lg p-6">
          <h2 className="text-2xl font-semibold mb-4 text-gray-800">Recent Transactions</h2>
          <TransactionTable />
        </div>
      </div>
    </main>
  );
}