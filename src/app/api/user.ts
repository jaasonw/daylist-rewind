import PocketBase from "pocketbase";

export async function POST(request: Request) {
  const pb = new PocketBase(process.env.BACKEND_URL);
  const data = await request.json();
  const authData = await pb
    .collection("users")
    .authWithPassword(
      process.env.ADMIN_USER ?? "",
      process.env.ADMIN_PASS ?? ""
    );

  console.log("Logged in successfully:", authData);

  return Response.json(data);
}
