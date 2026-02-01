import PlaidLink from '@/components/PlaidLink';
import TransactionTable from '@/components/TransactionTable';

export default function Home() {
  return (
    <main className="dashboard-container">
      {/* Dashboard Header */}
      <div className="dashboard-header">
        <div>
          <h1 className="header-title">
            Financial <span className="gradient-text">Agent</span>
          </h1>
          <p className="header-subtitle">AI-Powered Transaction Intelligence</p>
        </div>
        <PlaidLink />
      </div>

      {/* Main Content Grid */}
      <div className="dashboard-content">
        <div className="glass-panel transactions-panel">
          <div className="panel-header">
            <h2 className="panel-title">Recent Transactions</h2>
            <div className="live-indicator">
              Live Data Connection
            </div>
          </div>
          <TransactionTable />
        </div>
      </div>
    </main>
  );
}