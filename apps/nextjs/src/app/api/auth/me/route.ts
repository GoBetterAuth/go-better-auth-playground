import { NextResponse } from "next/server";

import { goBetterAuthClientServer } from "@/lib/gba-client-server";

export async function GET() {
  try {
    const response = await goBetterAuthClientServer.getMe();
    return NextResponse.json(response);
  } catch (error: any) {
    return NextResponse.json({ message: error?.message }, { status: 500 });
  }
}
