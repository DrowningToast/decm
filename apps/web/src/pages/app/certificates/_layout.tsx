import { Outlet } from "react-router-dom";

const CertificatesLayout = () => {
    return (
        <>
            {/* Background image */}
            <div className="absolute bottom-0 right-0 w-[424px] h-[424px] md:w-[500px] md:h-[500px] opacity-20 pointer-events-none">
                <img
                    src="/assets/passport.webp"
                    alt=""
                    className="w-full h-full object-cover object-center"
                />
            </div>
            <Outlet />
        </>
    );
};

export default CertificatesLayout;
