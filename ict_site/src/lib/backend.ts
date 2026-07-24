import packageJson from "../../package.json";

const WS_YEAR = new Date().getFullYear();

export const BE_POOL = process.env.BE_POOL ?? "http://ict_rest:36665";
export const WS_CONF = {
  name: process.env.WS_NAME ?? "OmniSight",
  version: packageJson.version,
  copyright: `© ${WS_YEAR}, ${process.env.WS_NAME ?? "OmniSight"}.`,
  meta: {
    title: process.env.WS_NAME ?? "OmniSight",
    description: process.env.WS_DESC ?? "OmniSight",
  },
};
