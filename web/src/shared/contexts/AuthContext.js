import { createContext } from "react";
import { NS } from "shared/utils";

export default createContext([{ isAuthenticated: true }, new NS("INIT"), () => {}, () => {}]);
