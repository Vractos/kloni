import { UserProvider } from "@auth0/nextjs-auth0/client";
import "./globals.css";
import type { Metadata } from "next";
import { Inter } from "next/font/google";
import NavBar from "./components/navbar";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Create Next App",
  description: "Generated by create next app",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="pt-BR" className="h-full bg-gray-100">
      <UserProvider>
        <body className={`${inter.className} h-full`}>
          <div className="min-h-full">
            <NavBar />
            <main>
              <div className="mx-auto max-w-7xl py-6 sm:px-6 lg:px-8">
              {children}
              </div>
            </main>
          </div>
        </body>
      </UserProvider>
    </html>
  );
}
