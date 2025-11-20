import Link from "next/link";

export default function Navbar() {
  return (
    <nav className="bg-black text-white p-4">
      <div className="container mx-auto flex gap-6">
        <Link href="/">Home</Link>
        <Link href="/transactions">Transactions</Link>
        <Link href="/rules">Rules</Link>
      </div>
    </nav>
  );
}