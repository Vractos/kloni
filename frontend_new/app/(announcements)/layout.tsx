import Search from "../components/search";

export default function AnnouncementPage({
  children,
  modal,
}: {
  children: React.ReactNode;
  modal: React.ReactNode;
}) {
  return (
    <>
      <div className="min-h-full bg-gray-100">
        <main className="max-h-screen" style={{ overflow: "auto" }}>
          <div className="mx-auto max-w-7xl py-6 sm:px-6 lg:px-8 ">
            <div className="mx-auto max-w-7xl py-6 px-4 sm:px-6 lg:px-8">
              <Search />
              {modal}
              {children}
            </div>
          </div>
        </main>
      </div>
    </>
  );
}
