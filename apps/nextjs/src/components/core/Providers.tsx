"use client";

import type { PropsWithChildren } from "react";

import { Toaster } from "sonner";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

const queryClient = new QueryClient();

export default function Providers({ children }: PropsWithChildren) {
  return (
    <>
      <QueryClientProvider client={queryClient}>
        <>{children}</>
        <Toaster />
      </QueryClientProvider>
    </>
  );
}
