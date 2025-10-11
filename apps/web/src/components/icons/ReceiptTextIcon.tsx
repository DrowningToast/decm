import * as React from "react";
import type { SVGProps } from "react";
const SvgReceiptTextIcon = (props: SVGProps<SVGSVGElement>) => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    fill="none"
    overflow="visible"
    preserveAspectRatio="none"
    style={{
      display: "block",
    }}
    viewBox="0 0 24 24"
    width="1em"
    height="1em"
    {...props}
  >
    <path
      stroke="var(--stroke-0, #E9DEDE)"
      strokeLinecap="round"
      strokeLinejoin="round"
      strokeWidth={2}
      d="M14 8H8m8 4H8m5 4H8M4 2v20l2-1 2 1 2-1 2 1 2-1 2 1 2-1 2 1V2l-2 1-2-1-2 1-2-1-2 1-2-1-2 1z"
    />
  </svg>
);
export default SvgReceiptTextIcon;
