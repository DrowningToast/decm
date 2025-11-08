import { useEffect } from "react";
import { useNavigate } from "@/router";

export default function IssuerHomePage() {
    const navigate = useNavigate();

    useEffect(() => {
        // Redirect to issuer sign page as the main issuer dashboard
        navigate("/issuer/sign");
    }, [navigate]);

    return null;
}
