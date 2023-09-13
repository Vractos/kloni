"use client"

import { UserCircleIcon, FingerPrintIcon, BellIcon, CubeIcon, CreditCardIcon, UsersIcon, KeyIcon, PuzzlePieceIcon } from '@heroicons/react/24/outline';
import Link from 'next/link';
import { usePathname } from 'next/navigation';

const secondaryNavigation = [
  { name: 'Geral', href: '/configuracoes', icon: UserCircleIcon},
  { name: 'Integração', href: '/configuracoes/integracao', icon: PuzzlePieceIcon},
  { name: 'Planos', href: '#', icon: CubeIcon},
  { name: 'Pagamento', href: '#', icon: CreditCardIcon},
]

function classNames(...classes: string[]) {
  return classes.filter(Boolean).join(' ')
}

export default function SideMenu() {
  const pathName = usePathname()

  return (
    <aside className="flex overflow-x-auto border-b border-gray-900/5 py-4 lg:block lg:w-64 lg:flex-none lg:border-0 lg:py-20">
      <nav className="flex-none px-4 sm:px-6 lg:px-0">
        <ul
          role="list"
          className="flex gap-x-3 gap-y-1 whitespace-nowrap lg:flex-col"
        >
          {secondaryNavigation.map((item) => (
            <li key={item.name}>
              <Link
                href={item.href}
                className={classNames(
                  pathName === item.href
                    ? "bg-gray-50 text-indigo-600"
                    : "text-gray-700 hover:text-indigo-600 hover:bg-gray-50",
                  "group flex gap-x-3 rounded-md py-2 pl-2 pr-3 text-sm leading-6 font-semibold"
                )}
              >
                <item.icon
                  className={classNames(
                    pathName === item.href
                      ? "text-indigo-600"
                      : "text-gray-400 group-hover:text-indigo-600",
                    "h-6 w-6 shrink-0"
                  )}
                  aria-hidden="true"
                />
                {item.name}
              </Link>
            </li>
          ))}
        </ul>
      </nav>
    </aside>
  );
}
