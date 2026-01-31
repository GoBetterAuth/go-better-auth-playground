import { NextRequest, NextResponse } from "next/server";

export function proxy(request: NextRequest) {
  // const cookieStore = await cookies();
  // const sessionCookie = cookieStore.get("gobetterauth.session_token");

  // const isAuthenticated = !!sessionCookie;
  // const pathname = request.nextUrl.pathname;

  // // Not authenticated → block dashboard
  // if (!isAuthenticated && pathname.startsWith("/dashboard")) {
  //   const url = request.nextUrl.clone();
  //   url.pathname = "/auth/sign-in";
  //   return NextResponse.redirect(url);
  // }

  // // Authenticated users should not see auth pages
  // if (
  //   isAuthenticated &&
  //   ["/auth/sign-in", "/auth/sign-up"].some((path) => pathname.startsWith(path))
  // ) {
  //   const url = request.nextUrl.clone();
  //   url.pathname = "/dashboard";
  //   return NextResponse.redirect(url);
  // }

  return NextResponse.next();
}

export const config = {
  matcher: [
    /*
     * Match all request paths except for the ones starting with:
     * - _next/static (static files)
     * - _next/image (image optimization files)
     * - favicon.ico (favicon file)
     * Feel free to modify this pattern to include more paths.
     */
    "/((?!_next/static|_next/image|favicon.ico|api|.*\\.(?:svg|png|jpg|jpeg|gif|webp)$).*)",
  ],
};
