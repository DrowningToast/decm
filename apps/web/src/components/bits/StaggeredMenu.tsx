import React, { useCallback, useLayoutEffect, useRef, useState } from "react";
import { gsap } from "gsap";
import { ThemisBlack, ThemisWhite } from "../icons";
import { cn } from "@/lib/utils";
import { AnimatePresence, motion } from "motion/react";
import { Typography } from "@/components/typography/typography";
import { useTranslation } from "react-i18next";

export interface StaggeredMenuItem {
    label: string;
    ariaLabel: string;
    link: string;
}
export interface StaggeredMenuSection {
    title: string;
    items: StaggeredMenuItem[];
}
export interface StaggeredMenuSocialItem {
    label: string;
    link: string;
}
export interface StaggeredMenuProps {
    position?: "left" | "right";
    colors?: string[];
    items?: StaggeredMenuItem[];
    sections?: StaggeredMenuSection[];
    socialItems?: StaggeredMenuSocialItem[];
    displaySocials?: boolean;
    displayItemNumbering?: boolean;
    className?: string;
    menuButtonColor?: string;
    openMenuButtonColor?: string;
    accentColor?: string;
    isFixed: boolean;
    changeMenuColorOnOpen?: boolean;
    onMenuOpen?: () => void;
    onMenuClose?: () => void;
}

export const StaggeredMenu: React.FC<StaggeredMenuProps> = ({
    position = "right",
    colors = ["var(--background-alt)", "var(--primary)"],
    items = [],
    sections = [],
    socialItems = [],
    displaySocials = true,
    displayItemNumbering = true,
    className,
    menuButtonColor = "var(--foreground-alt)",
    openMenuButtonColor = "var(--foreground)",
    changeMenuColorOnOpen = true,
    accentColor = "var(--primary)",
    isFixed = false,
    onMenuOpen,
    onMenuClose,
}: StaggeredMenuProps) => {
    const { t } = useTranslation();
    const [open, setOpen] = useState(false);
    const openRef = useRef(false);

    const panelRef = useRef<HTMLDivElement | null>(null);
    const preLayersRef = useRef<HTMLDivElement | null>(null);
    const preLayerElsRef = useRef<HTMLElement[]>([]);

    const plusHRef = useRef<HTMLSpanElement | null>(null);
    const plusVRef = useRef<HTMLSpanElement | null>(null);
    const iconRef = useRef<HTMLSpanElement | null>(null);

    const textInnerRef = useRef<HTMLSpanElement | null>(null);
    const textWrapRef = useRef<HTMLSpanElement | null>(null);

    const openTlRef = useRef<gsap.core.Timeline | null>(null);
    const closeTweenRef = useRef<gsap.core.Tween | null>(null);
    const spinTweenRef = useRef<gsap.core.Timeline | null>(null);
    const textCycleAnimRef = useRef<gsap.core.Tween | null>(null);
    const colorTweenRef = useRef<gsap.core.Tween | null>(null);

    const toggleBtnRef = useRef<HTMLButtonElement | null>(null);
    const busyRef = useRef(false);

    const itemEntranceTweenRef = useRef<gsap.core.Tween | null>(null);

    useLayoutEffect(() => {
        const ctx = gsap.context(() => {
            const panel = panelRef.current;
            const preContainer = preLayersRef.current;

            const plusH = plusHRef.current;
            const plusV = plusVRef.current;
            const icon = iconRef.current;

            console.log(panel, plusH, plusV, icon);
            if (!panel || !plusH || !plusV || !icon) {
                return;
            }

            let preLayers: HTMLElement[] = [];
            if (preContainer) {
                preLayers = Array.from(
                    preContainer.querySelectorAll(".sm-prelayer"),
                ) as HTMLElement[];
            }
            preLayerElsRef.current = preLayers;

            const offscreen = position === "left" ? -100 : 100;
            gsap.set([panel, ...preLayers], { xPercent: offscreen });

            gsap.set(plusH, { transformOrigin: "50% 50%", rotate: 0 });
            gsap.set(plusV, { transformOrigin: "50% 50%", rotate: 90 });
            gsap.set(icon, { rotate: 0, transformOrigin: "50% 50%" });

            if (toggleBtnRef.current) gsap.set(toggleBtnRef.current, { color: menuButtonColor });
        });
        return () => ctx.revert();
    }, [menuButtonColor, position]);

    const buildOpenTimeline = useCallback(() => {
        const panel = panelRef.current;
        const layers = preLayerElsRef.current;
        if (!panel) return null;

        openTlRef.current?.kill();
        if (closeTweenRef.current) {
            closeTweenRef.current.kill();
            closeTweenRef.current = null;
        }
        itemEntranceTweenRef.current?.kill();

        const itemEls = Array.from(panel.querySelectorAll(".sm-panel-itemLabel")) as HTMLElement[];
        const numberEls = Array.from(
            panel.querySelectorAll(".sm-panel-list[data-numbering] .sm-panel-item"),
        ) as HTMLElement[];
        const socialTitle = panel.querySelector(".sm-socials-title") as HTMLElement | null;
        const socialLinks = Array.from(panel.querySelectorAll(".sm-socials-link")) as HTMLElement[];

        const layerStates = layers.map((el) => ({
            el,
            start: Number(gsap.getProperty(el, "xPercent")),
        }));
        const panelStart = Number(gsap.getProperty(panel, "xPercent"));

        if (itemEls.length) gsap.set(itemEls, { yPercent: 140, rotate: 10 });
        if (numberEls.length) gsap.set(numberEls, { "--sm-num-opacity": 0 } as gsap.TweenVars);
        if (socialTitle) gsap.set(socialTitle, { opacity: 0 });
        if (socialLinks.length) gsap.set(socialLinks, { y: 25, opacity: 0 });

        const tl = gsap.timeline({ paused: true });

        layerStates.forEach((ls, i) => {
            tl.fromTo(
                ls.el,
                { xPercent: ls.start },
                { xPercent: 0, duration: 0.5, ease: "power4.out" },
                i * 0.07,
            );
        });

        const lastTime = layerStates.length ? (layerStates.length - 1) * 0.07 : 0;
        const panelInsertTime = lastTime + (layerStates.length ? 0.08 : 0);
        const panelDuration = 0.65;

        tl.fromTo(
            panel,
            { xPercent: panelStart },
            { xPercent: 0, duration: panelDuration, ease: "power4.out" },
            panelInsertTime,
        );

        if (itemEls.length) {
            const itemsStartRatio = 0.15;
            const itemsStart = panelInsertTime + panelDuration * itemsStartRatio;

            tl.to(
                itemEls,
                {
                    yPercent: 0,
                    rotate: 0,
                    duration: 1,
                    ease: "power4.out",
                    stagger: { each: 0.1, from: "start" },
                },
                itemsStart,
            );

            if (numberEls.length) {
                tl.to(
                    numberEls,
                    {
                        duration: 0.6,
                        ease: "power2.out",
                        "--sm-num-opacity": 1,
                        stagger: { each: 0.08, from: "start" },
                    } as gsap.TweenVars,
                    itemsStart + 0.1,
                );
            }
        }

        if (socialTitle || socialLinks.length) {
            const socialsStart = panelInsertTime + panelDuration * 0.4;

            if (socialTitle)
                tl.to(socialTitle, { opacity: 1, duration: 0.5, ease: "power2.out" }, socialsStart);
            if (socialLinks.length) {
                tl.to(
                    socialLinks,
                    {
                        y: 0,
                        opacity: 1,
                        duration: 0.55,
                        ease: "power3.out",
                        stagger: { each: 0.08, from: "start" },
                        onComplete: () => {
                            gsap.set(socialLinks, { clearProps: "opacity" });
                        },
                    },
                    socialsStart + 0.04,
                );
            }
        }

        openTlRef.current = tl;
        return tl;
    }, []);

    const playOpen = useCallback(() => {
        if (busyRef.current) return;
        busyRef.current = true;
        const tl = buildOpenTimeline();
        if (tl) {
            tl.eventCallback("onComplete", () => {
                busyRef.current = false;
            });
            tl.play(0);
        } else {
            busyRef.current = false;
        }
    }, [buildOpenTimeline]);

    const playClose = useCallback(() => {
        openTlRef.current?.kill();
        openTlRef.current = null;
        itemEntranceTweenRef.current?.kill();

        const panel = panelRef.current;
        const layers = preLayerElsRef.current;
        if (!panel) return;

        const all: HTMLElement[] = [...layers, panel];
        closeTweenRef.current?.kill();

        const offscreen = position === "left" ? -100 : 100;

        closeTweenRef.current = gsap.to(all, {
            xPercent: offscreen,
            duration: 0.32,
            ease: "power3.in",
            overwrite: "auto",
            onComplete: () => {
                const itemEls = Array.from(
                    panel.querySelectorAll(".sm-panel-itemLabel"),
                ) as HTMLElement[];
                if (itemEls.length) gsap.set(itemEls, { yPercent: 140, rotate: 10 });

                const numberEls = Array.from(
                    panel.querySelectorAll(".sm-panel-list[data-numbering] .sm-panel-item"),
                ) as HTMLElement[];
                if (numberEls.length)
                    gsap.set(numberEls, { "--sm-num-opacity": 0 } as gsap.TweenVars);

                const socialTitle = panel.querySelector(".sm-socials-title") as HTMLElement | null;
                const socialLinks = Array.from(
                    panel.querySelectorAll(".sm-socials-link"),
                ) as HTMLElement[];
                if (socialTitle) gsap.set(socialTitle, { opacity: 0 });
                if (socialLinks.length) gsap.set(socialLinks, { y: 25, opacity: 0 });

                busyRef.current = false;
            },
        });
    }, [position]);

    const animateIcon = useCallback((opening: boolean) => {
        const icon = iconRef.current;
        const h = plusHRef.current;
        const v = plusVRef.current;
        if (!icon || !h || !v) return;

        spinTweenRef.current?.kill();

        if (opening) {
            // ensure container never rotates
            gsap.set(icon, { rotate: 0, transformOrigin: "50% 50%" });
            spinTweenRef.current = gsap
                .timeline({ defaults: { ease: "power4.out" } })
                .to(h, { rotate: 45, duration: 0.5 }, 0)
                .to(v, { rotate: -45, duration: 0.5 }, 0);
        } else {
            spinTweenRef.current = gsap
                .timeline({ defaults: { ease: "power3.inOut" } })
                .to(h, { rotate: 0, duration: 0.35 }, 0)
                .to(v, { rotate: 90, duration: 0.35 }, 0)
                .to(icon, { rotate: 0, duration: 0.001 }, 0);
        }
    }, []);

    const animateColor = useCallback(
        (opening: boolean) => {
            const btn = toggleBtnRef.current;
            if (!btn) return;
            colorTweenRef.current?.kill();
            if (changeMenuColorOnOpen) {
                const targetColor = opening ? openMenuButtonColor : menuButtonColor;
                colorTweenRef.current = gsap.to(btn, {
                    color: targetColor,
                    delay: 0.18,
                    duration: 0.3,
                    ease: "power2.out",
                });
            } else {
                gsap.set(btn, { color: menuButtonColor });
            }
        },
        [openMenuButtonColor, menuButtonColor, changeMenuColorOnOpen],
    );

    React.useEffect(() => {
        if (toggleBtnRef.current) {
            if (changeMenuColorOnOpen) {
                const targetColor = openRef.current ? openMenuButtonColor : menuButtonColor;
                gsap.set(toggleBtnRef.current, { color: targetColor });
            } else {
                gsap.set(toggleBtnRef.current, { color: menuButtonColor });
            }
        }
    }, [changeMenuColorOnOpen, menuButtonColor, openMenuButtonColor]);

    const animateText = useCallback((opening: boolean) => {
        const inner = textInnerRef.current;
        if (!inner) return;

        textCycleAnimRef.current?.kill();

        const currentLabel = opening ? "Menu" : "Close";
        const targetLabel = opening ? "Close" : "Menu";
        const cycles = 3;

        const seq: string[] = [currentLabel];
        let last = currentLabel;
        for (let i = 0; i < cycles; i++) {
            last = last === "Menu" ? "Close" : "Menu";
            seq.push(last);
        }
        if (last !== targetLabel) seq.push(targetLabel);
        seq.push(targetLabel);

        // Render text lines dynamically
        if (inner) {
            inner.innerHTML = seq
                .map((line) => `<span class="sm-toggle-line">${line}</span>`)
                .join("");
        }
        gsap.set(inner, { yPercent: 0 });

        const lineCount = seq.length;
        const finalShift = ((lineCount - 1) / lineCount) * 100;

        textCycleAnimRef.current = gsap.to(inner, {
            yPercent: -finalShift,
            duration: 0.5 + lineCount * 0.07,
            ease: "power4.out",
        });
    }, []);

    const toggleMenu = useCallback(() => {
        const target = !openRef.current;
        openRef.current = target;
        setOpen(target);

        if (target) {
            onMenuOpen?.();
            playOpen();
        } else {
            onMenuClose?.();
            playClose();
        }

        animateIcon(target);
        animateColor(target);
        animateText(target);
    }, [playOpen, playClose, animateIcon, animateColor, animateText, onMenuOpen, onMenuClose]);

    return (
        <div
            className={`sm-scope z-40 ${isFixed ? "fixed top-0 left-0 w-screen h-screen overflow-hidden pointer-events-none" : "w-full h-full"}`}
        >
            <div
                className={
                    (className ? className + " " : "") +
                    "staggered-menu-wrapper relative w-full h-full z-40 pointer-events-none"
                }
                style={
                    accentColor
                        ? ({ ["--sm-accent"]: accentColor } as React.CSSProperties)
                        : undefined
                }
                data-position={position}
                data-open={open || undefined}
            >
                <div
                    ref={preLayersRef}
                    className={cn(
                        "sm-prelayers absolute top-0 right-0 bottom-0 pointer-events-none z-[5]",
                    )}
                    aria-hidden="true"
                >
                    {(() => {
                        const raw =
                            colors && colors.length
                                ? colors.slice(0, 4)
                                : ["var(--background-alt)", "var(--background)"];
                        const arr = [...raw];
                        if (arr.length >= 3) {
                            const mid = Math.floor(arr.length / 2);
                            arr.splice(mid, 1);
                        }
                        return arr.map((c, i) => (
                            <div
                                key={i}
                                className="sm-prelayer absolute top-0 right-0 h-full w-full translate-x-0"
                                style={{ background: c }}
                            />
                        ));
                    })()}
                </div>

                <header
                    className="staggered-menu-header absolute top-0 left-0 w-full flex items-center justify-between px-3 md:px-6 py-1.5 pointer-events-auto z-20 transition-all backdrop-blur-[2px] shadow-[0px_4px_16px_0px_rgba(0,0,0,0.1)] overflow-hidden bg-[rgba(217,217,217,0.1)]"
                    aria-label="Main navigation header"
                >
                    <div className="sm-logo flex items-center select-none" aria-label="Logo">
                        <div className="md:hidden">
                            <AnimatePresence mode="wait">
                                {!open ? (
                                    <motion.div
                                        initial={{ opacity: 0, y: 10 }}
                                        animate={{ opacity: 1, y: 0 }}
                                        exit={{ opacity: 0, y: 10 }}
                                        transition={{ duration: 0.3 }}
                                        layout
                                        key="open"
                                    >
                                        <ThemisWhite className="size-9 md:size-12" />
                                    </motion.div>
                                ) : (
                                    <motion.div
                                        initial={{ opacity: 0, y: 10 }}
                                        animate={{ opacity: 1, y: 0 }}
                                        exit={{ opacity: 0, y: 10 }}
                                        transition={{ duration: 0.3 }}
                                        layout
                                        key="closed"
                                    >
                                        <ThemisBlack className="size-9 md:size-12" />
                                    </motion.div>
                                )}
                            </AnimatePresence>
                        </div>
                        <div className="hidden md:block">
                            <ThemisWhite className="size-9 md:size-12" />
                        </div>
                    </div>

                    <button
                        ref={toggleBtnRef}
                        className={`sm-toggle relative inline-flex items-center gap-[0.3rem] bg-transparent border-0 cursor-pointer font-medium leading-none overflow-visible pointer-events-auto`}
                        aria-label={open ? "Close menu" : "Open menu"}
                        aria-expanded={open}
                        aria-controls="staggered-menu-panel"
                        onClick={toggleMenu}
                        type="button"
                    >
                        <span
                            ref={textWrapRef}
                            className="sm-toggle-textWrap relative inline-block h-[1em] overflow-hidden whitespace-nowrap w-[var(--sm-toggle-width,auto)] min-w-[var(--sm-toggle-width,auto)]"
                            aria-hidden="true"
                        ></span>
                        <motion.span
                            layout
                            ref={iconRef}
                            animate={{
                                color: open ? "var(--background)" : "var(--foreground)",
                            }}
                            className={cn(
                                "sm-icon relative w-[14px] h-[14px] shrink-0 inline-flex items-center justify-center [will-change:transform]",
                            )}
                            aria-hidden="true"
                        >
                            <span
                                ref={plusHRef}
                                className="sm-icon-line absolute left-1/2 top-1/2 w-full h-[2px] bg-current rounded-[2px] -translate-x-1/2 -translate-y-1/2 [will-change:transform]"
                            />
                            <span
                                ref={plusVRef}
                                className="sm-icon-line sm-icon-line-v absolute left-1/2 top-1/2 w-full h-[2px] bg-current rounded-[2px] -translate-x-1/2 -translate-y-1/2 [will-change:transform]"
                            />
                        </motion.span>
                    </button>
                </header>

                <aside
                    id="staggered-menu-panel"
                    ref={panelRef}
                    className="staggered-menu-panel absolute top-0 w-[clamp(260px,38vw,420px)] h-full bg-foreground flex flex-col p-[6em_2em_2em_2em] overflow-y-auto z-10 backdrop-blur-[12px] [-webkit-backdrop-filter:blur(12px)] right-0 [.staggered-menu-wrapper[data-position='left']_&]:right-auto [.staggered-menu-wrapper[data-position='left']_&]:left-0 pointer-events-auto"
                    aria-hidden={!open}
                >
                    <div className="sm-panel-inner flex-1 flex flex-col gap-5">
                        {sections && sections.length > 0 ? (
                            sections.map((section, sectionIdx) => (
                                <div
                                    key={section.title + sectionIdx}
                                    className="flex flex-col gap-3"
                                >
                                    <Typography
                                        variant="header"
                                        tag="h3"
                                        className="m-0 text-sm tracking-[0.06px] uppercase"
                                        color="background-alt"
                                    >
                                        {section.title}
                                    </Typography>
                                    <ul
                                        className="sm-panel-list list-none m-0 p-0 flex flex-col gap-2"
                                        role="list"
                                        data-numbering={displayItemNumbering || undefined}
                                    >
                                        {section.items.map((it, idx) => (
                                            <li
                                                className="relative overflow-hidden leading-none"
                                                key={it.label + idx}
                                            >
                                                <a
                                                    className="sm-panel-item relative cursor-pointer leading-none tracking-[-2px] uppercase transition-colors duration-150 ease-linear inline-block no-underline pr-[1.4em] hover:text-[var(--sm-accent,#ff0000)]"
                                                    href={it.link}
                                                    aria-label={it.ariaLabel}
                                                    data-index={sectionIdx * 100 + idx + 1}
                                                >
                                                    <Typography
                                                        variant="header"
                                                        tag="span"
                                                        className="sm-panel-itemLabel inline-block [transform-origin:50%_100%] will-change-transform text-[4rem] font-semibold text-black"
                                                    >
                                                        {it.label}
                                                    </Typography>
                                                </a>
                                            </li>
                                        ))}
                                    </ul>
                                </div>
                            ))
                        ) : (
                            <ul
                                className="sm-panel-list list-none m-0 p-0 flex flex-col gap-2"
                                role="list"
                                data-numbering={displayItemNumbering || undefined}
                            >
                                {items && items.length ? (
                                    items.map((it, idx) => (
                                        <li
                                            className="relative overflow-hidden leading-none"
                                            key={it.label + idx}
                                        >
                                            <a
                                                className="sm-panel-item relative text-black font-semibold text-[4rem] cursor-pointer leading-none tracking-[-2px] uppercase transition-colors duration-150 ease-linear inline-block no-underline pr-[1.4em] hover:text-[var(--sm-accent,#ff0000)]"
                                                href={it.link}
                                                aria-label={it.ariaLabel}
                                                data-index={idx + 1}
                                            >
                                                <span className="sm-panel-itemLabel inline-block [transform-origin:50%_100%] will-change-transform">
                                                    {it.label}
                                                </span>
                                            </a>
                                        </li>
                                    ))
                                ) : (
                                    <li
                                        className="relative overflow-hidden leading-none"
                                        aria-hidden="true"
                                    >
                                        <span className="sm-panel-item relative text-black font-semibold text-[4rem] cursor-pointer leading-none tracking-[-2px] uppercase transition-colors duration-150 ease-linear inline-block no-underline pr-[1.4em]">
                                            <span className="sm-panel-itemLabel inline-block [transform-origin:50%_100%] will-change-transform">
                                                {t("nav.noItems")}
                                            </span>
                                        </span>
                                    </li>
                                )}
                            </ul>
                        )}

                        {displaySocials && socialItems && socialItems.length > 0 && (
                            <div
                                className="sm-socials mt-auto pt-8 flex flex-col gap-3"
                                aria-label="Social links"
                            >
                                <Typography
                                    variant="header"
                                    tag="h3"
                                    className="m-0 text-base font-medium [color:var(--sm-accent,#ff0000)]"
                                >
                                    {t("nav.socials")}
                                </Typography>
                                <ul
                                    className="sm-socials-list list-none m-0 p-0 flex flex-row items-center gap-4 flex-wrap"
                                    role="list"
                                >
                                    {socialItems.map((s, i) => (
                                        <li key={s.label + i}>
                                            <a
                                                href={s.link}
                                                target="_blank"
                                                rel="noopener noreferrer"
                                                className="sm-socials-link text-[1.2rem] font-medium text-[#111] no-underline relative inline-block py-[2px] transition-[color,opacity] duration-300 ease-linear hover:text-[var(--sm-accent,#ff0000)] focus-visible:outline-2 focus-visible:outline-[var(--sm-accent,#ff0000)] focus-visible:outline-offset-[3px]"
                                            >
                                                {s.label}
                                            </a>
                                        </li>
                                    ))}
                                </ul>
                            </div>
                        )}
                    </div>
                </aside>
            </div>

            <style>{`
/* Core wrapper and positioning */
.sm-scope .staggered-menu-wrapper { position: relative; width: 100%; height: 100%; z-index: 40; }
.sm-scope .sm-logo { display: flex; align-items: center; user-select: none; }
.sm-scope .sm-logo-img { display: block; height: 32px; width: auto; object-fit: contain; }

/* Toggle button styles */
.sm-scope .sm-toggle { position: relative; display: inline-flex; align-items: center; gap: 0.3rem; background: transparent; border: none; cursor: pointer; color: #e9e9ef; font-weight: 500; line-height: 1; overflow: visible; }
.sm-scope .sm-toggle:focus-visible { outline: 2px solid #ffffffaa; outline-offset: 4px; border-radius: 4px; }
.sm-scope .sm-toggle-textWrap { position: relative; margin-right: 0.5em; display: inline-block; height: 1em; overflow: hidden; white-space: nowrap; width: var(--sm-toggle-width, auto); min-width: var(--sm-toggle-width, auto); }
.sm-scope .sm-toggle-textInner { display: flex; flex-direction: column; line-height: 1; }
.sm-scope .sm-toggle-line { display: block; height: 1em; line-height: 1; }

/* Icon animations */
.sm-scope .sm-icon { position: relative; width: 14px; height: 14px; flex: 0 0 14px; display: inline-flex; align-items: center; justify-content: center; will-change: transform; }
.sm-scope .sm-icon-line { position: absolute; left: 50%; top: 50%; width: 100%; height: 2px; background: currentColor; border-radius: 2px; transform: translate(-50%, -50%); will-change: transform; }
.sm-scope .sm-line { display: none !important; }

/* Pre-layers background animation */
.sm-scope .sm-prelayers { position: absolute; top: 0; right: 0; bottom: 0; width: clamp(260px, 38vw, 420px); pointer-events: none; z-index: 5; }
.sm-scope [data-position='left'] .sm-prelayers { right: auto; left: 0; }
.sm-scope .sm-prelayer { position: absolute; top: 0; right: 0; height: 100%; width: 100%; transform: translateX(0); }

/* Panel inner */
.sm-scope .sm-panel-inner { flex: 1; display: flex; flex-direction: column; gap: 1.25rem; }

/* Menu item label - required for GSAP animations */
.sm-scope .sm-panel-itemLabel { display: inline-block; will-change: transform; transform-origin: 50% 100%; }

/* Numbering system with CSS counters */
.sm-scope .sm-panel-list[data-numbering] { counter-reset: smItem; }
.sm-scope .sm-panel-list[data-numbering] .sm-panel-item::after { counter-increment: smItem; content: counter(smItem, decimal-leading-zero); position: absolute; top: 0; right: 0; font-size: 18px; font-weight: 400; color: var(--sm-accent, #ff0000); letter-spacing: 0; pointer-events: none; user-select: none; opacity: var(--sm-num-opacity, 0); transform: translateY(-20%); }

/* Social links hover effects - complex sibling selectors */
.sm-scope .sm-socials-list .sm-socials-link { opacity: 1; transition: opacity 0.3s ease; }
.sm-scope .sm-socials-list:hover .sm-socials-link:not(:hover) { opacity: 0.35; }
.sm-scope .sm-socials-list:focus-within .sm-socials-link:not(:focus-visible) { opacity: 0.35; }
.sm-scope .sm-socials-list .sm-socials-link:hover,
.sm-scope .sm-socials-list .sm-socials-link:focus-visible { opacity: 1; }

/* Responsive overrides */
@media (max-width: 1024px) { .sm-scope .staggered-menu-panel { width: 100%; left: 0; right: 0; } .sm-scope .staggered-menu-wrapper[data-open] .sm-logo-img { filter: invert(100%); } }
@media (max-width: 640px) { .sm-scope .staggered-menu-panel { width: 100%; left: 0; right: 0; } .sm-scope .staggered-menu-wrapper[data-open] .sm-logo-img { filter: invert(100%); } }
      `}</style>
        </div>
    );
};

export default StaggeredMenu;
