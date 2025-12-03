"use client";

import type { PropsWithChildren } from "react";

import { Toaster } from 'sonner';

export default function Providers({ children }: PropsWithChildren) {
  return (
    <>
      <>{children}</>
      <Toaster />
    </>
  );
}
