import { describe, it, expect } from "vitest";
import {
    toEventRegistrationConfigStatus,
    toEventRegistrationConfigStatusNumber,
} from "./event.utils";
import type { RequirementStatus } from "@/components/ui/requirement-item";

describe("event.utils", () => {
    describe("toEventRegistrationConfigStatus", () => {
        it("should return 'not_required' for undefined", () => {
            expect(toEventRegistrationConfigStatus(undefined)).toBe("not_required");
        });

        it("should return 'not_required' for 0", () => {
            expect(toEventRegistrationConfigStatus(0)).toBe("not_required");
        });

        it("should return 'required' for 1", () => {
            expect(toEventRegistrationConfigStatus(1)).toBe("required");
        });

        it("should return 'optional' for 2", () => {
            expect(toEventRegistrationConfigStatus(2)).toBe("optional");
        });

        it("should return 'not_required' for invalid numbers", () => {
            expect(toEventRegistrationConfigStatus(99)).toBe("not_required");
            expect(toEventRegistrationConfigStatus(-1)).toBe("not_required");
        });
    });

    describe("toEventRegistrationConfigStatusNumber", () => {
        it("should return 0 for undefined", () => {
            expect(toEventRegistrationConfigStatusNumber(undefined)).toBe(0);
        });

        it("should return 0 for 'not_required'", () => {
            expect(toEventRegistrationConfigStatusNumber("not_required")).toBe(0);
        });

        it("should return 1 for 'required'", () => {
            expect(toEventRegistrationConfigStatusNumber("required")).toBe(1);
        });

        it("should return 2 for 'optional'", () => {
            expect(toEventRegistrationConfigStatusNumber("optional")).toBe(2);
        });

        it("should return 0 for invalid status", () => {
            expect(toEventRegistrationConfigStatusNumber("invalid" as RequirementStatus)).toBe(0);
        });
    });

    describe("bidirectional conversion", () => {
        it("should correctly convert back and forth", () => {
            const testCases: Array<{ number: number; status: RequirementStatus }> = [
                { number: 0, status: "not_required" },
                { number: 1, status: "required" },
                { number: 2, status: "optional" },
            ];

            testCases.forEach(({ number, status }) => {
                expect(toEventRegistrationConfigStatus(number)).toBe(status);
                expect(toEventRegistrationConfigStatusNumber(status)).toBe(number);
            });
        });

        it("should maintain consistency in round-trip conversion", () => {
            const numbers = [0, 1, 2];

            numbers.forEach((num) => {
                const status = toEventRegistrationConfigStatus(num);
                const convertedBack = toEventRegistrationConfigStatusNumber(status);
                expect(convertedBack).toBe(num);
            });
        });
    });
});
