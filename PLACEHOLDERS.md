# PLACEHOLDERS — things needing Evan's hands / decisions

- [ ] **License choice**: repo ships with MIT (permissive default). Say the word to switch
      (Apache-2.0, GPL, proprietary, etc.).
- [ ] **GPL/LGPL tooling decision**: mtkclient (GPL-3.0) and Heimdall (LGPL) are the
      standard open-source backends for MediaTek and Samsung flashing. They are NOT
      bundled or linked — `doctor` only detects them and prints install hints.
      Decide whether Dro1dRa1n may shell out to them as separate processes.
- [ ] **GitHub push**: commits are local-only until you say push.
- [ ] **Real-device testing**: automated tests use fixtures; on-device verification
      (Samsung test-mode FRP flow, Xiaomi sideload, fastboot reset) needs a physical
      device in your hands.
- [ ] **Windows smoke test**: `dro1drain.exe` cross-compiles clean; run it on a real
      Windows box with platform-tools + Samsung USB drivers to confirm.
- [ ] **Driver installs**: Windows needs the Samsung USB driver / Google USB driver;
      `doctor` prints where to get them but cannot install them for you.
