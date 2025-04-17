import CaddyEditor from "@/app/CaddyEditor";
import {Suspense} from "react";
import HistoryViewer from "@/app/HistoryViewer";

export default function Home() {
  return (
    <div className="w-full min-h-[100svh]">
      <HistoryViewer/>
      <br/>
      <CaddyEditor />
    </div>
  );
}
