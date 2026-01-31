import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

import { TooltipProvider } from "~/components/ui/tooltip";
import { Toaster } from "~/components/ui/sonner";
import type { PropsWithChildren } from "react";

const queryClient = new QueryClient();

export default function Providers({ children }: PropsWithChildren) {
  return (
    <QueryClientProvider client={queryClient}>
      <TooltipProvider>
        <>{children}</>
        <Toaster />
      </TooltipProvider>
    </QueryClientProvider>
  );
}
