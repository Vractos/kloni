import { UserProvider } from "@auth0/nextjs-auth0/client";
import "./globals.css";
import type { Metadata } from "next";
import { Inter } from "next/font/google";
import React from "react";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Kloni",
  icons: {
    other: {
      rel: "icon",
      url: "data:image/svg+xml,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 100 100'><text y='.9em' font-size='90'>🪨</text></svg>",
    },
  },
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
