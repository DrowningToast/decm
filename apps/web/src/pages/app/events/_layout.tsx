import { Outlet } from "react-router-dom";

const EventsLayout = () => {
    return (
        <>
            {/* Background image */}
            <div className="absolute bottom-0 right-8 w-[129px] aspect-auto h-auto md:w-[320px] opacity-30 pointer-events-none">
                <img
                    src="/assets/hand.webp"
                    alt=""
                    className="w-full h-full object-cover object-center"
                />
            </div>
            <Outlet />
        </>
    );
};

export default EventsLayout;
