import { UserProvider } from "@auth0/nextjs-auth0/client";
import "./globals.css";
import type { Metadata } from "next";
import { Inter } from "next/font/google";
import React from "react";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Kloni",
};
export default function RootLayout({
  children,

}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="pt-BR" className="h-full">
      <UserProvider>
        <body
          className={`${inter.className} h-full`}
          style={{ overflow: "hidden" }}
          >
          {children}
        </body>
      </UserProvider>
    </html>
  );
}
