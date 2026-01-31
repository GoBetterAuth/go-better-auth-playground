import { useEffect } from "react";
import type { UseQueryResult } from "@tanstack/react-query";
import { useQuery } from "@tanstack/react-query";
import type { GetMeResponse } from "go-better-auth";

import { toast } from "~/hooks/use-toast";
import { goBetterAuthClient } from "~/lib/gba-client";

export async function fetchMe(): Promise<GetMeResponse> {
  const data = await goBetterAuthClient.getMe<GetMeResponse>();
  return data;
}

export function useMe(): UseQueryResult<GetMeResponse, any> {
  const query = useQuery({
    queryKey: ["me"],
    queryFn: fetchMe,
    retry: false,
    staleTime: 1000 * 60,
  });

  useEffect(() => {
    if (query.error) {
      // Lightweight feedback for unexpected failures
      toast({
        title: "Failed to load profile",
        description: query.error?.message ?? String(query.error),
      });
    }
  }, [query.error]);

  return query;
}
