import { useEffect } from "react";
import { useNavigate } from "react-router-dom";

export default function EmptyHostPage() {
    const navigate = useNavigate();

    useEffect(() => {
        navigate("/host/home");
    }, [navigate]);

    return null;
}
