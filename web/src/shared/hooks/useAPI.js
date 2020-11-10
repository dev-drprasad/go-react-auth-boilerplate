import { message } from "antd";
import { useCallback, useContext, useEffect, useMemo, useState } from "react";
import { AuthContext } from "shared/contexts";
// import userService from "src/shared/services/user.service";
import { NS } from "shared/utils";

export const cr = "\n";
export const tab = "\t";

function getType(data) {
  if (data === null) return "Null";
  if (data === undefined) return "Undefined";
  if (typeof data === "string") return "String";
  if (typeof data === "number" && !Number.isNaN(data)) return "Number";
  if (Number.isNaN(data)) return "NaN";
  if (typeof data === "boolean") return "Boolean";
  if (data instanceof Array) return "Array"; // always should be before `Object`
  if (data instanceof Object) return "Object";

  return "";
}

// https://stackoverflow.com/a/2117523
function uuidv4() {
  return ([1e7] + -1e3 + -4e3 + -8e3 + -1e11).replace(/[018]/g, (c) =>
    (c ^ (crypto.getRandomValues(new Uint8Array(1))[0] & (15 >> (c / 4)))).toString(16)
  );
}

const defaultFetchOptions = {
  headers: {
    Accept: "application/json",
    "Content-Type": "application/json",
  },
};

export function useFetch(url, opts) {
  const [[rId, fresh], setParams] = useState(() => ["", false]);
  const [response, setResponse] = useState([undefined, new NS("INIT")]);

  const refresh = useCallback(() => setParams([uuidv4(), true]), []);

  useEffect(() => {
    if (url && opts) {
      setParams([uuidv4(), false]);
    }
  }, [url, opts]);

  useEffect(() => {
    if (!url || !rId) return;

    const abortctrl = new AbortController();

    setResponse(([, s]) => [undefined, s.clone("LOADING")]);
    const startTime = performance.now();

    // recursive merge might be better solution
    const finalopts = {
      ...defaultFetchOptions,
      ...opts,
      headers: {
        ...defaultFetchOptions.headers,
        ...opts?.headers,
        ...(fresh && { "X-Clear-Cache": true }),
        "X-Request-ID": rId,
      },
    };

    if (finalopts.headers["Content-Type"] === null) {
      delete finalopts.headers["Content-Type"];
    }
    fetch(url, finalopts)
      .then(async (res) => {
        if (abortctrl.signal.aborted) return;

        const responseTime = performance.now() - startTime;
        const cached = !!res.headers.get("X-Browser-Cache");
        let body;

        try {
          body = await res.json();
        } catch (e) {
          const message = "Invalid JSON response from API";
          console.error(`${cr}API Error:${cr}${tab}URL: ${url}${cr}${tab}Msg: ${message}${cr}${tab}Code: ${res.status}`);
          setResponse(([, s]) => [undefined, s.clone("ERROR", "", res.status, responseTime, rId, cached)]);
          return;
        }

        if (res.status >= 400) {
          const errorType = body.error || "";
          const isInternalError = !errorType || errorType === "Internal Server Error";
          const message = !isInternalError ? body.message || "" : "";
          setResponse(([, s]) => [undefined, s.clone("ERROR", message, res.status, responseTime, rId, cached, false)]);
          return;
        }

        const dataType = getType(body);
        const hasData = dataType !== "Null" && (dataType === "Array" ? body.length > 0 : true);

        setResponse(([, s]) => [body, s.clone("SUCCESS", "", res.status, responseTime, rId, cached, hasData)]);
      })
      .catch((err) => {
        if (abortctrl.signal.aborted) return;
        const responseTime = performance.now() - startTime;
        console.error(`${cr}API Error:${cr}${tab}URL: ${url}${cr}${tab}Msg: ${err.message}${cr}${tab}Code: 0`);
        setResponse(([, s]) => [undefined, s.clone("ERROR", "", 0, responseTime, rId)]);
      });

    return () => abortctrl.abort();
  }, [fresh, rId, url]); //dont add `opts`, adding will cause two renders when error happens

  return [response[0], response[1], refresh];
}

const API_BASE_URL = process.env.REACT_APP_API_BASE_URL || window.location.origin;
export default function useAPI(urlpath, extraOptions) {
  const [user, , logout] = useContext(AuthContext);
  const url = urlpath && new URL(urlpath, API_BASE_URL).toString();

  const options = useMemo(
    () => ({
      ...extraOptions,
      headers: {
        Authorization: `Bearer ${user?.token || ""}`,
        ...extraOptions?.headers,
      },
    }),
    [extraOptions]
  ); // Don't add `user` . Adding will cause re-render when
  // token change (ex: /login api)

  const [data, status, refresh] = useFetch(url, options);

  useEffect(() => {
    if (status.statusCode === 401) {
      logout();
    }
    if (status.isError && !status.errorCaught && status.statusCode >= 400) {
      message.error("Oops! Something went wrong.", 3);
    }
  }, [status, logout]);

  return [data, status, refresh];
}
