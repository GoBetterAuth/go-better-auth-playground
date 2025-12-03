import type { Metadata } from "next";
import './globals.css';

import { Raleway } from 'next/font/google';

import Providers from '@/components/core/Providers';
import { cn } from '@/lib/utils';

const primaryFont = Raleway({
  variable: "--font-primary",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "GoBetterAuth Playground",
  description:
    "An example modern authentication solution built with Go and Next.js.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body
        className={cn(`${primaryFont.variable} antialiased`)}
        suppressHydrationWarning
      >
        <Providers>{children}</Providers>
      </body>
    </html>
  );
}
