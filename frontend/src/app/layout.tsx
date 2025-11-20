import "../styles/globals.css";
import Navbar from "../components/Navbar";

export const metadata = { title: "Autonomous Transaction Categorizer" };

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body>
        <Navbar />
        <main className="container mx-auto p-6">{children}</main>
      </body>
    </html>
  );
}