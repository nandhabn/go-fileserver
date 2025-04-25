import type { Route } from "./+types/home";
import { Welcome } from "../welcome/welcome";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "One Piece - Streamer" },
    { name: "description", content: "One Piece - Streamer" },
  ];
}

export default function Home() {
  return <Welcome />;
}
