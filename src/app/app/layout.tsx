import { AuthProvider } from "@/components/PocketBaseAuthProvider";
import { ThemeProvider } from "@/components/ThemeProvider";

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <ThemeProvider
      attribute="class"
      defaultTheme="system"
      enableSystem
    >
      <AuthProvider>{children}</AuthProvider>
    </ThemeProvider>
  );
}
