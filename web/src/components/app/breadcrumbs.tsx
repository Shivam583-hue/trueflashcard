import Link from "next/link";
import { CaretRightIcon } from "@phosphor-icons/react/dist/ssr";

type Crumb = { label: string; href?: string };

export function Breadcrumbs({ items }: { items: Crumb[] }) {
  return (
    <nav className="flex items-center gap-1.5 text-sm">
      {items.map((item, i) => {
        const last = i === items.length - 1;
        return (
          <span key={i} className="flex items-center gap-1.5">
            {item.href && !last ? (
              <Link
                href={item.href}
                className="text-neutral-500 transition-colors duration-150 hover:text-neutral-200"
              >
                {item.label}
              </Link>
            ) : (
              <span
                className={last ? "font-medium text-neutral-100" : "text-neutral-500"}
              >
                {item.label}
              </span>
            )}
            {!last && <CaretRightIcon size={13} className="text-neutral-700" />}
          </span>
        );
      })}
    </nav>
  );
}
