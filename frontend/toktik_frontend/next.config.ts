// import type { NextConfig } from "next";

// const nextConfig: NextConfig = {
//   /* config options here */
// };

// export default nextConfig;


import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: "standalone", // ✅ Required for Docker deployment
  reactStrictMode: true, // Optional but recommended
};

export default nextConfig;

