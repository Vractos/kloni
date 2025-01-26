import Search from "@/components/search";
import { Suspense } from "react";

export default function AnnouncementPage({
  children,
  modal,
}: {
  children: React.ReactNode;
  modal: React.ReactNode;
}) {
  return (
    <>
      {modal}
      {children}
    </>
  );
}
